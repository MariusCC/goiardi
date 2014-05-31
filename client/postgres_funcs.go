/*
 * Copyright (c) 2013-2014, Jeremy Bingham (<jbingham@gmail.com>)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package client

import (
	"github.com/ctdk/goiardi/data_store"
	"github.com/ctdk/goiardi/util"
	"database/sql"
	"log"
	"net/http"
	"strings"
)

func getClientPostgreSQL(name string) (*Client, error) {
	client := new(Client)
	stmt, err := data_store.Dbh.Prepare("select c.name, nodename, validator, admin, o.name, public_key, certificate FROM goiardi.clients c JOIN goiardi.organizations o on c.organization_id = o.id WHERE c.name = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(name)
	err = client.fillClientFromSQL(row)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) savePostgreSQL() util.Gerror {
	tx, err := data_store.Dbh.Begin()
	if err != nil {
		gerr := util.CastErr(err)
		return gerr
	}
	_, err = tx.Exec("SELECT goiardi.merge_clients($1, $2, $3, $4, $5, $6)", c.Name, c.NodeName, c.Validator, c.Admin, c.pubKey, c.Certificate);
	if err != nil {
		tx.Rollback()
		gerr := util.CastErr(err)
		if strings.HasPrefix(err.Error(), "a user with") {
			gerr.SetStatus(http.StatusConflict)
		}
		return gerr
	}
	tx.Commit()
	return nil
}

func (c *Client) deletePostgreSQL() error {
	tx, err := data_store.Dbh.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM goiardi.clients WHERE name = $1", c.Name)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *Client) renamePostgreSQL(new_name string) util.Gerror {
	tx, err := data_store.Dbh.Begin()
	if err != nil {
		gerr := util.Errorf(err.Error())
		return gerr
	}
	_, err = tx.Exec("SELECT goiardi.rename_client($1, $2)", c.Name, new_name)
	if err != nil {
		tx.Rollback()
		gerr := util.Errorf(err.Error())
		if strings.HasPrefix(err.Error(), "a user with") || strings.Contains(err.Error(), "already exists, cannot rename") {
			gerr.SetStatus(http.StatusConflict)
		} else {
			gerr.SetStatus(http.StatusInternalServerError)
		}
		return gerr
	}
	tx.Commit()
	return nil
}

func numAdminsPostgreSQL() int {
	var numAdmins int
	stmt, err := data_store.Dbh.Prepare("SELECT count(*) FROM goiardi.clients WHERE admin = TRUE")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&numAdmins)
	if err != nil {
		log.Fatal(err)
	}
	return numAdmins
}

func getListPostgreSQL() []string {
	var client_list []string
	rows, err := data_store.Dbh.Query("SELECT name FROM goiardi.clients")
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		rows.Close()
		return client_list
	}
	client_list = make([]string, 0)
	for rows.Next() {
		var client_name string
		err = rows.Scan(&client_name)
		if err != nil {
			log.Fatal(err)
		}
		client_list = append(client_list, client_name)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return client_list
}