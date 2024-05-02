// mongostore contains gorilla session store.
package mongostore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libsession "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewStore creates new mongo store.
//
// The client argument is the mongo db client where sessions will be saved.
//
// Keys are defined in pairs to allow key rotation, but the common case is
// to set a single authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for
// encryption. The encryption key can be set to nil or omitted in the last
// pair, but the authentication key is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes.
// The encryption key, if set, must be either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256 modes.
func NewStore(client mongo.DbClient, keyPairs ...[]byte) *MongoStore {
	return &MongoStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:     "/",
			MaxAge:   int(30 * 24 * time.Hour.Seconds()),
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		},
		collection: client.Collection(mongo.SessionMongoCollection),
	}
}

// MongoStore stores sessions in mongo db.
type MongoStore struct {
	Codecs           []securecookie.Codec
	Options          *sessions.Options // default configuration
	collection       mongo.DbCollection
	autoCleanRunning bool
}

// sessionData represents mongo collection structure.
type sessionData struct {
	ID      primitive.ObjectID `bson:"_id"`
	Data    string             `bson:"data"`
	Expires int64              `bson:"expires"`
}

// Get returns a session for the given name after adding it to the registry.
//
// It returns a new session if the sessions doesn't exist. Access IsNew on
// the session to check if it is an existing session or a new one.
//
// It returns a new session and an error if the session exists but could
// not be decoded.
func (s *MongoStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// The difference between New() and Get() is that calling New() twice will
// decode the session data twice, while Get() registers and reuses the same
// decoded session after the first call.
func (s *MongoStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		if err == nil {
			err = s.load(r.Context(), session)
			session.IsNew = false
		}
		var securecookieError securecookie.Error
		if errors.As(err, &securecookieError) {
			// if securecookie decode failed (for example due changed key), then it's a new session
			err = nil
		}
	}
	return session, err
}

// Save adds a single session to the response.
//
// If the Options.MaxAge of the session is <= 0 then the session file will be
// deleted from the store path. With this process it enforces the properly
// session cookie handling so no need to trust in the cookie management in the
// web browser.
func (s *MongoStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Delete if max-age is <= 0
	if session.Options.MaxAge <= 0 {
		if err := s.erase(r.Context(), session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	d := time.Duration(session.Options.MaxAge) * time.Second
	expires := time.Now().Add(d)

	if err := s.save(r.Context(), session, expires); err != nil {
		return err
	}

	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
	if err != nil {
		return err
	}

	cookie := sessions.NewCookie(session.Name(), encoded, session.Options)
	cookie.Expires = expires
	http.SetCookie(w, cookie)

	return nil
}

func (s *MongoStore) StartAutoClean(ctx context.Context, timeout time.Duration) {
	if s.autoCleanRunning {
		return
	}
	s.autoCleanRunning = true
	ticker := time.NewTicker(timeout)
	defer func() {
		ticker.Stop()
		s.autoCleanRunning = false
	}()
	for {
		select {
		case <-ticker.C:
			_ = s.clean(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *MongoStore) GetActiveSessionsCount(ctx context.Context) (int64, error) {
	return s.collection.CountDocuments(ctx, bson.M{
		"expires": bson.M{"$gt": time.Now().Unix()},
	})
}

// load finds session in mongo collection and decodes its content into session.Values.
func (s *MongoStore) load(ctx context.Context, session *sessions.Session) error {
	if session.ID == "" {
		return nil
	}

	id, err := primitive.ObjectIDFromHex(session.ID)
	if err != nil {
		return err
	}

	cursor, err := s.collection.Find(ctx, bson.M{
		"_id":     id,
		"expires": bson.M{"$gt": time.Now().Unix()},
	})

	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var data sessionData
		err := cursor.Decode(&data)

		if err != nil {
			return err
		}

		if err = securecookie.DecodeMulti(session.Name(), data.Data,
			&session.Values, s.Codecs...); err != nil {
			return err
		}

		return nil
	}

	return libsession.ErrNoSession
}

// erase deletes session from mongo collection.
func (s *MongoStore) erase(ctx context.Context, session *sessions.Session) error {
	if session.ID == "" {
		return nil
	}

	id, err := primitive.ObjectIDFromHex(session.ID)
	if err != nil {
		return err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (s *MongoStore) ExpireSessions(ctx context.Context, user string, provider string) error {
	_, err := s.collection.DeleteMany(ctx, bson.M{"user": user, "provider": provider})
	return err
}

func (s *MongoStore) ExpireSessionsByUserIDs(ctx context.Context, ids []string) error {
	_, err := s.collection.DeleteMany(ctx, bson.M{"user": bson.M{"$in": ids}})
	return err
}

// save writes encoded session.Values to mongo collection.
func (s *MongoStore) save(ctx context.Context, session *sessions.Session, expires time.Time) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values, s.Codecs...)
	if err != nil {
		return err
	}

	doc := bson.M{
		"data":    encoded,
		"expires": expires.Unix(),
	}

	provider, ok := session.Values["provider"]
	if ok {
		doc["provider"] = provider
	}

	userID, ok := session.Values["user"]
	if ok {
		doc["user"] = userID
	}

	if session.ID == "" {
		res, err := s.collection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}

		if id, ok := res.(primitive.ObjectID); ok {
			session.ID = id.Hex()

			return nil
		}

		return fmt.Errorf("cannot parse string from %v", res)
	}

	id, err := primitive.ObjectIDFromHex(session.ID)
	if err != nil {
		return err
	}

	res, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": doc})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("fail to modify session : session %v not found", session.ID)
	}

	return nil
}

// clean deletes all expired sessions from mongo collection.
func (s *MongoStore) clean(ctx context.Context) error {
	_, err := s.collection.DeleteMany(ctx, bson.M{"expires": bson.M{"$lte": time.Now().Unix()}})

	return err
}
