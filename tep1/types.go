/*
 * Copyright (C) 2019 The TesraSupernet Authors
 * This file is part of The TesraSupernet library.
 *
 * The TesraSupernet is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The TesraSupernet is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The TesraSupernet.  If not, see <http://www.gnu.org/licenses/>.
 */
package tep1

import (
	"fmt"
	"github.com/TesraSupernet/Tesra/common"
	"math/big"
)

type State struct {
	From   common.Address
	To     common.Address
	Amount *big.Int
}

type Tep1TransferEvent struct {
	Name   string
	From   common.Address
	To     common.Address
	Amount *big.Int
}

func (this *Tep1TransferEvent) String() string {
	return fmt.Sprintf("name %s, from %s, to %s, amount %s", this.Name, this.From.ToBase58(), this.To.ToBase58(),
		this.Amount.String())
}
