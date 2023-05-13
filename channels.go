package wildfrostbot

import bolt "go.etcd.io/bbolt"

const channelBucket string = "channels"

// GetChannelPerm returns true if the given channel has message commands enabled
func (h *DiscordHandler) GetChannelPerm(id string) bool {
	var res bool
	h.ChannelDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(channelBucket))
		v := b.Get([]byte(id))
		if v != nil {
			res = true
		}
		return nil
	})
	return res
}

// SetChannelPerm sets the permission for message commands in a given channel
func (h *DiscordHandler) SetChannelPerm(id string, val bool) error {
	if val {
		return h.ChannelDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(channelBucket))
			return b.Put([]byte(id), []byte("true"))
		})
	} else {
		return h.ChannelDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(channelBucket))
			return b.Delete([]byte(id))
		})
	}
}

// CreateDB creates a bbolt store for persisting data about permissions
func (h *DiscordHandler) CreateDB(path string) error {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}

	h.ChannelDB = db

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(channelBucket))
		return err
	})
}
