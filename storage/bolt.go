package storage

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"go.etcd.io/bbolt"
)

var ErrNoteNotFound = errors.New("笔记不存在")

type Note struct {
	ID        uint64    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

var bucketName = []byte("notes")

type Store struct {
	db *bbolt.DB
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) AddNote(content string) (uint64, error) {
	var id uint64
	err := s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		id, _ = b.NextSequence()
		note := Note{
			ID:        id,
			Content:   content,
			CreatedAt: time.Now(),
		}
		data, err := json.Marshal(note)
		if err != nil {
			return err
		}
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, id)
		return b.Put(key, data)
	})
	return id, err
}

func (s *Store) ListNotes() ([]Note, error) {
	var notes []Note
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.ForEach(func(k, v []byte) error {
			var note Note
			if err := json.Unmarshal(v, &note); err != nil {
				return err
			}
			notes = append(notes, note)
			return nil
		})
	})
	return notes, err
}

func (s *Store) GetNote(id uint64) (*Note, error) {
	var note *Note
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, id)
		v := b.Get(key)
		if v == nil {
			return ErrNoteNotFound
		}
		var n Note
		if err := json.Unmarshal(v, &n); err != nil {
			return err
		}
		note = &n
		return nil
	})
	return note, err
}

func (s *Store) DeleteNote(id uint64) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, id)
		v := b.Get(key)
		if v == nil {
			return ErrNoteNotFound
		}
		return b.Delete(key)
	})
}

func (s *Store) SearchNotes(keyword string) ([]Note, error) {
	var notes []Note
	keywordLower := strings.ToLower(keyword)
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.ForEach(func(k, v []byte) error {
			var note Note
			if err := json.Unmarshal(v, &note); err != nil {
				return err
			}
			if strings.Contains(strings.ToLower(note.Content), keywordLower) {
				notes = append(notes, note)
			}
			return nil
		})
	})
	return notes, err
}
