/*
* Copyright (c) 2014 Santiago Arias | Remy Jourde
*
* Permission to use, copy, modify, and distribute this software for any
* purpose with or without fee is hereby granted, provided that the above
* copyright notice and this permission notice appear in all copies.
*
* THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
* WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
* MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
* ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
* WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
* ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
* OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package helpers

import "fmt"

var (
	themes = []string{"frogideas", "bythepool", "heatwave", "summerwarmth"}
)

// UserImageURL returns the user image URL for a given name and id.
//
func UserImageURL(name string, id int64) string {
	t := themes[id%4]
	return fmt.Sprintf("https://www.tinygraphs.com/spaceinvaders/%s?theme=%s&numcolors=%d", name, t, (id%2)+2)
}

// TournamentImageURL returns the tournament image URL for a given name and id.
//
func TournamentImageURL(name string, id int64) string {
	t := themes[id%4]
	return fmt.Sprintf("https://www.tinygraphs.com/labs/isogrids/hexa/%s?theme=%s&numcolors=%d", name, t, (id%2)+2)
}

// TeamImageURL returns the team image URL for a given name and id.
//
func TeamImageURL(name string, id int64) string {
	t := themes[id%4]
	return fmt.Sprintf("https://www.tinygraphs.com/labs/isogrids/hexa16/%s?theme=%s&numcolors=%d", name, t, (id%2)+2)
}
