package storage

import (
	"math/rand"
	"testing"
	"time"
)

var coms = []Comment{
	{User_id: 5, Text: "папуасы нервно курят", PubTime: 1698997341, ParentType: "A", ParentID: 2347},
	{User_id: 23, Text: "потом узнаем, что будет тогда", PubTime: 1699193457, ParentType: "A", ParentID: 34},
	{User_id: 4, Text: "ну и к чему это все", PubTime: 1699393031, ParentType: "C", ParentID: 5},
	{User_id: 5, Text: "сам такой", PubTime: 1699413037, ParentType: "A", ParentID: 23},
	{User_id: 2, Text: "а вот когда я был маленький, себе такого не позволял", PubTime: 1699513038, ParentType: "C", ParentID: 67},

	{User_id: 4, Text: "упал под стол, пополз умирать под плинтус", PubTime: 1699543033, ParentType: "A", ParentID: 87},
	{User_id: 3, Text: "с языка снял", PubTime: 1699843033, ParentType: "C", ParentID: 12},
	{User_id: 765, Text: "гаф гаф пизда", PubTime: 1699853037, ParentType: "A", ParentID: 3},
	{User_id: 5, Text: "абырвалг", PubTime: 1699953007, ParentType: "A", ParentID: 1},
	{User_id: 222, Text: "ой девочки, что будет", PubTime: 1699957007, ParentType: "A", ParentID: 46},
	{User_id: 23, Text: "кто здесь?!", PubTime: 1699997007, ParentType: "A", ParentID: 7},
}

func TestNew(t *testing.T) {
	_, err := New("@localhost:5432/comments")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDB_Comment(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	c := Comment{

		Text:       "Test comm3",
		PubTime:    1699219670,
		ParentType: "A",
		ParentID:   453,
	}
	db, err := New("@localhost:5432/comments")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.SaveComment(c)
	if err != nil {
		t.Fatal(err)
	}
	comments, err := db.CommentList("A", 453)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", comments)
}

func TestDB_SaveComment(t *testing.T) {
	type args struct {
		c Comment
	}
	dbp, err := New("@localhost:5432/comments")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		db      *DB
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "com1",
			db:      dbp,
			args:    args{c: coms[0]},
			want:    0,
			wantErr: false,
		},
		{
			name:    "com2",
			db:      dbp,
			args:    args{c: coms[1]},
			want:    0,
			wantErr: false,
		},
		{
			name:    "com3",
			db:      dbp,
			args:    args{c: coms[2]},
			want:    0,
			wantErr: false,
		},
		{
			name:    "com4",
			db:      dbp,
			args:    args{c: coms[3]},
			want:    0,
			wantErr: false,
		},
		{
			name:    "com5",
			db:      dbp,
			args:    args{c: coms[4]},
			want:    0,
			wantErr: false,
		},
		{
			name:    "com6",
			db:      dbp,
			args:    args{c: coms[5]},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.SaveComment(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.SaveComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.want {
				t.Errorf("DB.SaveComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_DeleteComment(t *testing.T) {
	type args struct {
		id int
	}
	dbp, err := New("@localhost:5432/comments")
	if err != nil {
		t.Fatal(err)
	}
	id_c, err := dbp.SaveComment(coms[8])
	if err != nil {
		t.Fatal(err, id_c)
	}
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		db      *DB
		args    args
		wantErr bool
	}{
		{
			name:    "com1",
			db:      dbp,
			args:    args{id: id_c},
			wantErr: false,
		},
		{
			name:    "com2",
			db:      dbp,
			args:    args{id: -1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.DeleteComment(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DB.DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
