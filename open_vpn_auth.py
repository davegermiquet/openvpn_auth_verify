#!/bin/python
import sqlite3
import argparse
import sys
import bcrypt
import os


def create_database(payload,database):

    with database:
        database.execute("CREATE TABLE USER(user VARCHAR UNIQUE, password BLOB)")
    return 1

def login_user(payload,database):

    if "username" in os.environ and "password" in os.environ:
        user = os.environ['username'] 
        password = os.environ['password']
        cur = database.cursor()
        cur.execute("SELECT password FROM USER where USER = ?",(user,))
        rows = cur.fetchall()
        cur.close()
        if len(rows) == 1:
            if bcrypt.checkpw(password.encode("utf-8"),rows[0][0]):
                print("passed")
                return 0
    return 1


def add_user_with_password(payload,database):

    hash_password = bcrypt.hashpw(payload['password'].encode("utf-8"), bcrypt.gensalt())
    with database:
        database.execute("INSERT INTO USER(user,password) VALUES(?,?)", (payload['user'],hash_password))
    return 1

def delete_user(payload,database):

    with database:
        database.execute("DELETE FROM USER WHERE user like '%"+payload['user'] + "%'")
    return 1

def invalid_command(payload,database):

    print("Invalid arguments")
    return -1

def parse_args():

    parser = argparse.ArgumentParser(description='Python Authentication Plugin For OpenVPN')
    parser.add_argument("-a","--add",help="Add user to database",   action='store_true')
    parser.add_argument("-d","--delete",help="Delete user to database",   action='store_true')
    parser.add_argument("-u","--user",help="User for Action")
    parser.add_argument("-p","--password",help="Password")
    parser.add_argument("-c","--createdb",help="create database",action="store_true")
    args = parser.parse_args()
    payload = None
    if args.createdb:
        command = create_database
    else:
        if args.add:
            if not args.user or not args.password:
                command = invalid_command
            else:
                payload = {
                    "user" : args.user,
                    "password": args.password
                }
                command = add_user_with_password


        elif args.delete:
            if args.user:
                command = delete_user
                payload = {
                    "user": args.user
                }
            else:
                command = invalid_command
        
        else:
            command = login_user 

    return {
    "command" : command,
    "payload": payload
    }

def main():
    DEFAULTDB = "user.db"
    con = sqlite3.connect(DEFAULTDB)
    try:
        command = parse_args() 
        ret_value = command['command'](payload=command['payload'],database = con)
    except Exception as ex:
        con.close()
        print(ex)
        sys.exit(-1)
    finally:
        con.close()

    sys.exit(ret_value)
if __name__ == "__main__":
    main()


