package session

import (
	"database/sql"
	_ "embed"
	"fmt"
	"main/core"
	"os"
	"path"
	"time"

	"github.com/flytam/filenamify"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

//go:embed schema.sql
var schema string

type session struct {
	listeners []core.Listener
	db        *sql.DB
	log       *logrus.Logger
	filepath  string
}

func New(log *logrus.Logger, md core.Metadata) (core.Session, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	filename, err := filenamify.Filenamify(fmt.Sprintf("Cook %s.sql", md.Start.Format("Jan _2 2006")), filenamify.Options{})
	if err != nil {
		return nil, err
	}

	sess := &session{
		listeners: make([]core.Listener, 0),
		log:       log,
		filepath:  path.Join(homedir, filename),
	}
	err = sess.SetMetadata(md)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func Open(log *logrus.Logger, filepath string) (core.Session, error) {
	sess := &session{
		listeners: make([]core.Listener, 0),
		log:       log,
		filepath:  filepath,
	}
	err := sess.open()
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *session) NewReading(r core.Reading) {
	err := s.open()
	if err != nil {
		s.log.Error(err)
		return
	}

	_, err = s.db.Exec("INSERT INTO readings (received, temp1, temp2) VALUES ($1, $2, $3)", r.Received.Unix(), r.Temp1, r.Temp2)
	if err != nil {
		s.log.Error(err)
	}

	for _, l := range s.listeners {
		l(s, r)
	}
}

func (s *session) GetReadings() ([]core.Reading, error) {
	err := s.open()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query("SELECT received, temp1, temp2 FROM readings ORDER BY received")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	readings := make([]core.Reading, 0)
	for rows.Next() {
		var reading core.Reading
		var receivedInt int64
		err = rows.Scan(
			&receivedInt,
			&reading.Temp1,
			&reading.Temp2,
		)
		if err != nil {
			return nil, err
		}

		reading.Received = time.Unix(receivedInt, 0)
		readings = append(readings, reading)
	}
	return readings, nil
}

func (s *session) SetMetadata(m core.Metadata) error {
	err := s.open()
	if err != nil {
		return err
	}

	keyValues := map[string]string{
		"food":   m.Food,
		"method": m.Method,
		"start":  m.Start.Format(time.RFC3339Nano),
	}
	for key, value := range keyValues {
		row := s.db.QueryRow("SELECT key FROM metadata WHERE key = $1", key)
		var err error
		if row.Scan() == sql.ErrNoRows {
			_, err = s.db.Exec("INSERT INTO metadata (key, value) VALUES ($1, $2)", key, value)
		} else {
			_, err = s.db.Exec("UPDATE metadata SET value = $1 WHERE key = $2", value, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *session) GetMetadata() (core.Metadata, error) {
	err := s.open()
	if err != nil {
		return core.Metadata{}, err
	}

	rows, err := s.db.Query("SELECT key, value FROM metadata")
	if err != nil {
		return core.Metadata{}, err
	}
	defer rows.Close()

	metadata := core.Metadata{}

	for rows.Next() {
		var key, value string
		err = rows.Scan(&key, &value)
		if err != nil {
			return core.Metadata{}, err
		}
		switch key {
		case "food":
			metadata.Food = value
		case "method":
			metadata.Method = value
		case "start":
			t, err := time.Parse(time.RFC3339Nano, value)
			if err != nil {
				return core.Metadata{}, err
			}
			metadata.Start = t
		}
	}

	return metadata, nil
}

func (s *session) AddListener(l core.Listener) {
	s.listeners = append(s.listeners, l)
}

func (s *session) open() error {
	if s.db != nil {
		return nil
	}

	_, err := os.Stat(s.filepath)
	create := os.IsNotExist(err)

	db, err := sql.Open("sqlite3", s.filepath)
	if err != nil {
		return err
	}

	if create {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}

	s.db = db

	return nil
}