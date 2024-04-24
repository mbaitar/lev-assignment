package view

import (
	"context"
	"fmt"
	"github.com/mbaitar/levenue-assignment/types"
	"strconv"
	"strings"
)

func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user, ok := ctx.Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}
	return user
}

func String(i int) string {
	return strconv.Itoa(i)
}

func Float64(i float64) string {
	return strconv.FormatFloat(i, 'f', 2, 64)
}

func FormatCurrency(i float64) string {
	formatted := fmt.Sprintf("%.2f", i)
	parts := strings.Split(formatted, ".")
	var integerPart []string
	integerStr := parts[0]

	for i, r := range integerStr {
		if i > 0 && (len(integerStr)-i)%3 == 0 {
			integerPart = append(integerPart, ",")
		}
		integerPart = append(integerPart, string(r))
	}
	integerFormatted := strings.Join(integerPart, "")
	return fmt.Sprintf("€%s.%s", integerFormatted, parts[1])
}
