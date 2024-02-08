#!/bin/python
import sqlite3
import argparse
import sys
import bcrypt
import os


def create_database(payload,database):

    with database:
        database.execute("CREATE TABLE USER(user VARCHAR UNIQUE, password BLOB)")
    
def login_user(payload,database):
    if "username" in os.environ and "password" in os.environ:
        user = os.environ['username'] 
        password = os.environ['password']
        cur = database.cursor()
        cur.execute("SELECT password FROM USER where USER = ?",(user,))
        rows = cur.fetchall()
        if len(rows) == 1:
            if bcrypt.checkpw(password.encode("utf-8"),rows[0][0]):
                print("passed")
                sys.exit(0)
    sys.exit(1)


def add_user_with_password(payload,database):
    hash_password = bcrypt.hashpw(payload['password'].encode("utf-8"), bcrypt.gensalt())
    with database:
        database.execute("INSERT INTO USER(user,password) VALUES(?,?)", (payload['user'],hash_password))
    
def delete_user(payload,database):
    with database:
        database.execute("DELETE FROM USER WHERE user like '%"+payload['user'] + "%'")

def invalid_command(payload,database):
    print("Invalid arguments")

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
        command['command'](payload=command['payload'],database = con)
    finally:
        con.close()

if __name__ == "__main__":
    main()


