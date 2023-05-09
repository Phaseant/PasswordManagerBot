package repository

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DataRepository struct {
	db        *mongo.Client
	secretKey string
}

// Add docs to DB
func (r DataRepository) Add(service, username, password, userID string) error {
	collection := r.db.Database(DBName).Collection(DataCol)
	doc := Document{
		Service:  service,
		Username: username,
		Password: r.createHash(password),
		UserID:   userID,
	}

	_, err := collection.InsertOne(context.Background(), doc)
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("service %s for username %s already exists", service, username)
	}
	if err != nil {
		return err
	}
	return nil
}

// Delete docs from DB
func (r DataRepository) Delete(serive, username, userID string) error {
	collection := r.db.Database(DBName).Collection(DataCol)
	filter := bson.D{{Key: "service", Value: serive}, {Key: "userID", Value: userID}, {Key: "username", Value: username}}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

// Get docs from DB
func (r DataRepository) Get(service, userID string) ([]ServiceCreds, error) {
	collection := r.db.Database(DBName).Collection(DataCol)

	var docs []Document

	records := []ServiceCreds{}

	filter := bson.D{{Key: "service", Value: service}, {Key: "userID", Value: userID}}
	cursor, _ := collection.Find(context.Background(), filter)
	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, err
	}
	if len(docs) > 0 {
		for _, doc := range docs {
			records = append(records, ServiceCreds{Login: doc.Username, Password: r.getPassword(doc.Password)})
		}
		return records, nil
	} else {
		return nil, errors.New("no documents found for this service")
	}
}

// Hash all passwords, secret key is in config
func (r DataRepository) createHash(password string) string {
	encrypted, err := encrypt(r.secretKey, password)
	if err != nil {
		log.Fatal(err)
	}
	return encrypted

}

// decrypt password from DB and return as string
func (r DataRepository) getPassword(password string) string {
	decrypted, err := decrypt(r.secretKey, password)
	if err != nil {
		log.Fatal(err)
	}
	return decrypted
}

// encrypt and decrypt functions
// Shoutout to https://gist.github.com/mickelsonm/e1bf365a149f3fe59119
func encrypt(keyStr string, message string) (encoded string, err error) {
	key := []byte(keyStr)
	//Create byte array from the input string
	plainText := []byte(message)

	//Create a new AES cipher using the key
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText), err
}

func decrypt(keyStr string, secure string) (decoded string, err error) {
	key := []byte(keyStr)
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	//IF DecodeString failed, exit:
	if err != nil {
		return
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
}
