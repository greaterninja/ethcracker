package main

import ( 
    "flag"
    "os"
    "bufio"
//    "fmt"
    "strings"
//    "encoding/json"
    "io/ioutil"
    "./crypto"
)

var templates [][]string 


// note, that variables are pointers
var pk = flag.String("pk", "", "Private key file")
var t = flag.String("t", "", "Pattern file")
var min_len = flag.Int("min_len", 8, "Minimum password length")
var max_len = flag.Int("max_len", 20, "Maximum password length")


func main() {
    flag.Parse()
    
    println( "------------------------------------------------")
    println( "Ethereum Password Cracker")
    println( "Author: @AlexNa ")
    println( "------------------------------------------------")
    println( "Private Key File:", *pk )
    println( "Template File:", *t )
    println( "Minimum password length", *min_len )
    println( "Maximum password length", *max_len )
    
    keyFileContent, err := ioutil.ReadFile( *pk )
    if err != nil { panic( err ) }

//    println( "Private key file content:", string( keyFileContent ) )

    key_version := 3
    key_v3, err_v3 := crypto.LoadKeyVersion3( keyFileContent )
    key_v1, err_v1 := crypto.LoadKeyVersion1( keyFileContent )

    if err_v3 != nil { 
        key_version = 1
        if err_v1 != nil { panic( err ) }
        //println( "Private key JSON (version 1):", string( pk_log ) )
    } else {
        //println( "Private key JSON (version 3):", string( pk_log ) )
    }

    println( "Key file version:", key_version)
    
    
    templates = make( [][]string, 0 )
    f, err := os.Open( *t )
    if err != nil { panic( err ) }

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        tl := strings.Split( scanner.Text(), " " )
        if len( tl ) >= 0 { templates = append( templates, tl ) }
    }

    if err := scanner.Err(); err != nil { panic( err ) }

    f.Close()
    
    println( "Template lines:", len( templates ) )
    
    n := 1
    for _,l := range templates { n *= len( l ) + 1 }
    n = n - 1 //exclude empty string
    println( "Total combinations:", n )
    

    //main cycle
    indexes := make( []int, len( templates ) )

    v := 0
    main: for {
        for i := 0; i < len( indexes ); i++ {
            
            if indexes[i] < len( templates[i] ) {
                indexes[i] = indexes[i] + 1
                break;
            } else {
                indexes[i] = 0
                if i == len( templates ) - 1 { break main }
            }
        } 
        
        v++
        s := ""
        for i := 0; i < len( indexes ); i++ {
            if indexes[i] > 0 { s = s + templates[i][indexes[i] - 1 ] }
        }
        
        if len(s) < *min_len { continue }
        if len(s) > *max_len { continue }
        
        println( v, "/", n, " ", s )
        
        if key_version == 3 {
            err = crypto.Test_pass_v3( key_v3, s )
        } else {
            err = crypto.Test_pass_v1( key_v1, s )
        }
        
        if err == nil {
            println( "Your password:", s )            
            return
        }
        
        
    }
           
           
    println( ":-( Sorry... password not found")           
} 