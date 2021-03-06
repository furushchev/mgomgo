package mgomgo

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"strings"
	"time"
	"io"
)

type DBParams struct {
	Host       string
	Database   string
	Collection string
	UserName   string
	Password   string
}

func NewDBParamsFromURI(uri string) (*DBParams, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "mongodb" {
		return nil, fmt.Errorf("mongodb:// scheme is only supported")
	}
	path := strings.Split(u.Path, "/")
	if len(path) == 2 {
		return nil, fmt.Errorf("No collection name specified")
	} else if len(path) != 3 {
		return nil, fmt.Errorf("database name and collection name must be specified")
	}

	var username = ""
	var password = ""
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}
	return &DBParams{
		Host:       u.Host,
		Database:   path[1],
		Collection: path[2],
		UserName:   username,
		Password:   password,
	}, nil
}


func Migrate(from, to string, conn int, timeout time.Duration) error {
	fromParams, err := NewDBParamsFromURI(from)
	if err != nil {
		return err
	}
	toParams, err := NewDBParamsFromURI(to)
	if err != nil {
		return err
	}

	logrus.Infof("connect as source to: %s\n", fromParams.Host)
	fromSession, err := mgo.DialWithTimeout(fromParams.Host, timeout)
	if err != nil {
		return err
	}
	defer fromSession.Clone()
	fromSession.SetMode(mgo.Monotonic, true)

	logrus.Infof("connect as destination to: %s\n", toParams.Host)
	toSession, err := mgo.DialWithTimeout(toParams.Host, timeout)
	if err != nil {
		return err
	}
	defer toSession.Close()
	toSession.SetMode(mgo.Monotonic, true)

	fromDB := fromSession.DB(fromParams.Database)
	if fromParams.UserName != "" {
		if err := fromDB.Login(fromParams.UserName, fromParams.Password); err != nil {
			return err
		}
	}
	fromC := fromDB.C(fromParams.Collection)
	logrus.Infof("connected source to %s.%s\n", fromC.Database.Name, fromC.Name)

	logrus.Infof("generating %d routines...", conn)
	iter := fromC.Find(bson.M{}).Iter()
	datachan := make(chan bson.M, conn)
	infochan := make(chan string, conn)
	errchan := make(chan error, conn)
	for i := 0; i < conn; i++ {
		go func(rnum int, s *mgo.Session, p *DBParams) {
			copySession := s.Copy()
			defer copySession.Close()
			copySession.SetMode(mgo.Monotonic, true)
			toDB := copySession.DB(p.Database)
			if toParams.UserName != "" {
				if err := toDB.Login(p.UserName, p.Password); err != nil {
					errchan <- err
				}
			}
			c := toDB.C(p.Collection)
			logrus.Infof("%d: connected to %s.%s\n", rnum, c.Database.Name, c.Name)
			for {
				rdata, ok := <- datachan
				if ok {
					if err := c.Insert(rdata); err != nil {
						if mgo.IsDup(err) {
							if oid, ok := rdata["_id"].(bson.ObjectId); ok {
								infochan <- fmt.Sprintf("%d: skipped %s", rnum, oid.Hex())
							}
							continue
						}
						if err == io.EOF {
							copySession.Refresh()
							continue
						}
						errchan <- err
						return
					} else {
						if oid, ok := rdata["_id"].(bson.ObjectId); ok {
							infochan <- fmt.Sprintf("%d: migrated %s", rnum, oid.Hex())
						} else {
							infochan <- fmt.Sprintf("%d: warn: cannot get _id: %v", rnum, rdata)
						}
					}
				}
			}
		}(i, toSession, toParams)
	}

	go func() {
		for {
			select {
			case info := <- infochan:
				logrus.Infoln(info)
			case err := <- errchan:
				logrus.Errorln(err)
			}
		}
	}()
	for {
		var data bson.M
		if ok := iter.Next(&data); !ok {
			logrus.Fatalf("no next data")
			close(datachan)
		} else {
			datachan <- data
		}
	}
	return nil
}
