package crypto

import ( 
    "encoding/json"
    //"io/ioutil"
    "github.com/pborman/uuid"
	"errors"
	"bytes"
	"encoding/hex"    
)

func LoadKeyVersion1( fileContent []byte ) ( *encryptedKeyJSONV1, error ) {
    key := new( encryptedKeyJSONV1)
    err := json.Unmarshal(fileContent, key )
    return key, err
}

func LoadKeyVersion3( fileContent []byte ) ( *encryptedKeyJSONV3, error ) {
    key := new( encryptedKeyJSONV3 )
    err := json.Unmarshal(fileContent, key )
    return key, err
}

func Test_pass_v1( k *encryptedKeyJSONV1, auth string ) error {
    _, _, err := MydecryptKeyV1(k, auth)
    return err
}

func Test_pass_v3( k *encryptedKeyJSONV3, auth string ) error {
    _, _, err := decryptKeyV3(k, auth)
    return err
}

func MydecryptKeyV1(keyProtected *encryptedKeyJSONV1, auth string) (keyBytes []byte, keyId []byte, err error) {
	keyId = uuid.Parse(keyProtected.Id)
	mac, err := hex.DecodeString(keyProtected.Crypto.MAC)
	if err != nil {
		return nil, nil, err
	}

	iv, err := hex.DecodeString(keyProtected.Crypto.CipherParams.IV)
	if err != nil {
		return nil, nil, err
	}

	cipherText, err := hex.DecodeString(keyProtected.Crypto.CipherText)
	if err != nil {
		return nil, nil, err
	}

	derivedKey, err := getKDFKey(keyProtected.Crypto, auth)
	if err != nil {
		return nil, nil, err
	}

	calculatedMAC := Sha3(derivedKey[16:32], cipherText)
	if !bytes.Equal(calculatedMAC, mac) {
		return nil, nil, errors.New("Decryption failed: MAC mismatch")
	}

	plainText, err := aesCBCDecrypt(Sha3(derivedKey[:16])[:16], cipherText, iv)
	if err != nil {
		return nil, nil, err
	}
	return plainText, keyId, err
}

