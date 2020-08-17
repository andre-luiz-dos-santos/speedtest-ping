package main

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func getInfoForAddr(addr string) (string, error) {
	info := "IP: " + addr + "\n"

	rows, err := ixcDB.Query(`
		SELECT
			radacct.username AS login,
			radacct.acctstarttime AS start,
			COALESCE( cliente_condominio.condominio, '' ) AS condominio,
			COALESCE( IF( radusuarios.endereco_padrao_cliente = 'N', radusuarios.bloco, cliente.bloco ), '' ) AS bloco,
			COALESCE( IF( radusuarios.endereco_padrao_cliente = 'N', radusuarios.apartamento, cliente.apartamento ), '' ) AS apartamento
		FROM
			radacct
		LEFT JOIN
			radusuarios ON ( radusuarios.login = radacct.username )
		LEFT JOIN
			cliente ON ( cliente.id = radusuarios.id_cliente )
		LEFT JOIN
			cliente_condominio ON ( cliente_condominio.id =
			IF( radusuarios.endereco_padrao_cliente = 'N', radusuarios.id_condominio, cliente.id_condominio ) )
		WHERE
			radacct.acctstoptime IS NULL
		AND
			radacct.framedipaddress = ?`, addr)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			login       string
			start       mysql.NullTime
			condominio  string
			bloco       string
			apartamento string
		)
		err = rows.Scan(&login, &start, &condominio, &bloco, &apartamento)
		if err != nil {
			return "", err
		}

		info += "--\n"
		info += fmt.Sprintf("Login: %v\n", login)
		info += fmt.Sprintf("Condomínio: %v\n", condominio)
		info += fmt.Sprintf("Bloco: %v\n", bloco)
		info += fmt.Sprintf("Apartamento: %v\n", apartamento)

		if start.Valid {
			info += fmt.Sprintf("Connection Start: %v\n",
				start.Time.Format("02/01/2006 15:04:05"))
			// IXC stores times in local time without time zone.
			// start.Time.Local().Format("02/01/2006 15:04:05 -0700"))
		}
	}

	return info, nil
}
