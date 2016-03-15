package crypto

import ( 
    "encoding/json"
    "os"
    //"io/ioutil"
    "github.com/pborman/uuid"
	"errors"
	"bytes"
	"encoding/hex"    
    "io/ioutil"    
)


type CrackerParams struct {
    key_version int 
    key_v1 *encryptedKeyJSONV1 
    key_v3 *encryptedKeyJSONV3 
    n int
}

func LoadKeyFile( params *CrackerParams, path string ) error {
    keyFileContent, err := ioutil.ReadFile( path )
    if err != nil { return err }

//    println( "Private key file content:", string( keyFileContent ) )

    params.key_version = 3
    params.key_v3, err = LoadKeyVersion3( keyFileContent )

    if err != nil { 
        params.key_version = 1
        params.key_v1, err = LoadKeyVersion1( keyFileContent )
        if err != nil { return err }
        //println( "Private key JSON (version 1):", string( pk_log ) )
    } else {
        //println( "Private key JSON (version 3):", string( pk_log ) )
    }

    println( "Key file version:", params.key_version)    
    return nil
}

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

func Test_pass( params *CrackerParams, s string ) error {

    var err error 
    
    if params.key_version == 3 {
        err = Test_pass_v3( params.key_v3, s )
    } else {
        err = Test_pass_v1( params.key_v1, s )
    }
    
    params.n++
    
    println( params.n, ": ", s )
        
    if err == nil {
        println( "" )            
        println( "" )            
        println( "-------------------------------------------------------------------------" )            
        println( "              CONGRATULATION !!! WE FOUND YOUR PASSWORD !!!" )            
        println( "-------------------------------------------------------------------------" )            
        println( "" )            
        println( "                  Password:", s )            
        println( "" )            
        println( "" )            
        println( "          Do not forget to donate some ETH to the developer:" )            
        println( "     Ethereum Address: 0x281694Fabfdd9735e01bB59942B18c469b6e3df6" )            
        println( "-------------------------------------------------------------------------" )            
        println( "" )            
        println( "" )            
        os.Exit( 0 );
    }
    
    return err
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

