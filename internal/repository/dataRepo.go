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

func (r *DataRepository) Add(service, username, password, userID string) error {
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
func (r *DataRepository) Delete(serive, userID string) error {
	collection := r.db.Database(DBName).Collection(DataCol)
	filter := bson.D{{Key: "service", Value: serive}, {Key: "userID", Value: userID}}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *DataRepository) Get(service, userID string) (string, string, error) {
	collection := r.db.Database(DBName).Collection(DataCol)

	var doc Document

	filter := bson.D{{Key: "service", Value: service}, {Key: "userID", Value: userID}}
	err := collection.FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		return "", "", err
	}
	log.Println("db.Get: returned login and password for user: ", userID, " and service: ", service)
	return doc.Username, r.getPassword(doc.Password), nil
}

func (r DataRepository) createHash(password string) string {
	encrypted, err := encrypt(r.secretKey, password)
	if err != nil {
		log.Fatal(err)
	}
	return encrypted

}

func (r DataRepository) getPassword(password string) string {
	decrypted, err := decrypt(r.secretKey, password)
	if err != nil {
		log.Fatal(err)
	}
	return decrypted
}

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
