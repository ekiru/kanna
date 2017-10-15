// The accounts package defines request handlers and database helpers
// for retrieving and working with people's account data.
package accounts

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"

	"golang.org/x/crypto/scrypt"

	"github.com/ekiru/kanna/actors"
	"github.com/ekiru/kanna/db"
)

// accounts.Model represents an account on this server.
type Model struct {
	// Username is used by the owner of the account to log-in to
	// this server.
	Username string
	// PasswordHash is a hash of the password for the account.
	PasswordHash []byte
	// PasswordHashVersion identifies which password hash algorithm
	// was used to encode PasswordHash.
	PasswordHashVersion PasswordHashAlgorithm
	// Actor is the main Actor belonging to the account. The
	// account may have permission to view Activities delivered to
	// other Actors or to author Activities as other Actors, but
	// this Actor represents this account specifically.
	Actor *actors.Model
}

// FromRow fills a Model with the data from a row returned by a
// database query from the Accounts table joined with the Actors table.
func (m *Model) FromRow(rows *sql.Rows) error {
	m.Actor = &actors.Model{}
	actor := m.Actor.Scanners()
	return rows.Scan(
		&m.Username,
		&m.PasswordHash,
		&m.PasswordHashVersion,
		actor["id"],
		actor["type"],
		actor["name"],
		actor["inbox"],
		actor["outbox"],
	)
}

// ByUsername retrieves a accounts.Model for the account with the
// supplied username, as well as the account's actor.
func ByUsername(ctx context.Context, username string) (*Model, error) {
	var model Model
	rows, err := db.DB(ctx).QueryContext(ctx,
		"select acct.username, acct.passwordHash, acct.passwordHashVersion, "+
			"acct.actorId, act.type, act.name, act.inbox, act.outbox "+
			"from Accounts acct join Actors act on acct.actorId = act.id "+
			"where username = ?",
		username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	if err = model.FromRow(rows); err != nil {
		return nil, err
	}
	return &model, nil
}

// Authenticate attempts to authenticate as an account.
func Authenticate(ctx context.Context, username string, password string) (*Model, error) {
	// TODO maybe avoid the user enum.
	user, err := ByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if !user.PasswordHashVersion.Matches(password, user.PasswordHash) {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

// A PasswordHashAlgorithm identifies a particular algorithm and set of
// parameters for password hashing to allow easily upgrading to slower
// or otherwise more secure algorithms or parameters in the future.
type PasswordHashAlgorithm int

const (
	// HashScrypt uses scrypt with N = 2^15, r = 8, p = 1, based on
	// the parameters in https://blog.filippo.io/the-scrypt-parameters/
	HashScrypt PasswordHashAlgorithm = iota
)

func (alg PasswordHashAlgorithm) Hash(password string, salt []byte) []byte {
	switch alg {
	case HashScrypt:
		if salt == nil {
			salt = make([]byte, 8)
			if _, err := rand.Read(salt); err != nil {
				panic(err)
			}
		}
		hash, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
		if err != nil {
			panic(err)
		}
		return append(hash, salt...)
	default:
		panic("Unrecognized hashing algorithm.")
	}
}

func (alg PasswordHashAlgorithm) Matches(password string, target []byte) bool {
	var hash []byte
	switch alg {
	case HashScrypt:
		salt := target[32:]
		hash = alg.Hash(password, salt)
	default:
		panic("Unrecognized hashing algorithm.")
	}

	// we could do a constant time compare but it doesn't matter here since we're comparing password hashes.
	return bytes.Equal(hash, target)
}
