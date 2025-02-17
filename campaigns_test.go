package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var camp = CampaignAd{
	Placeholder: "placholder",
	Ratio:       0.5,
	Id:          "id",
	Probability: 1,
	Fallback:    true,
	Ad: Ad{
		Source:      "source",
		Image:       "image",
		Link:        "http://link.com",
		Description: "desc",
		Company:     "company",
	},
}

func TestAddAndFetchCampaigns(t *testing.T) {
	migrateDatabase()
	initializeDatabase()
	defer tearDatabase()
	defer dropDatabase()

	err := addCampaign(context.Background(), ScheduledCampaignAd{
		CampaignAd: camp,
		Start:      time.Now().Add(time.Hour * -1),
		End:        time.Now().Add(time.Hour),
	})
	assert.Nil(t, err)
	err = addOrUpdateUserTags(context.Background(), "1", []string{"javascript"})
	assert.Nil(t, err)

	var res []CampaignAd
	res, err = fetchCampaigns(context.Background(), time.Now(), "1")
	assert.Nil(t, err)
	assert.Equal(t, []CampaignAd{camp}, res)
}

func TestFetchExpiredCampaigns(t *testing.T) {
	migrateDatabase()
	initializeDatabase()
	defer tearDatabase()
	defer dropDatabase()

	err := addCampaign(context.Background(), ScheduledCampaignAd{
		CampaignAd: camp,
		Start:      time.Now().Add(time.Hour * -2),
		End:        time.Now().Add(time.Hour * -1),
	})
	assert.Nil(t, err)

	var res []CampaignAd
	res, err = fetchCampaigns(context.Background(), time.Now(), "1")
	assert.Nil(t, err)
	assert.Equal(t, []CampaignAd(nil), res)
}

func TestFetchCampaignsWithTags(t *testing.T) {
	migrateDatabase()
	initializeDatabase()
	defer tearDatabase()
	defer dropDatabase()

	err := addCampaign(context.Background(), ScheduledCampaignAd{
		CampaignAd: camp,
		Start:      time.Now().Add(time.Hour * -1),
		End:        time.Now().Add(time.Hour),
	})
	assert.Nil(t, err)
	err = addCampaign(context.Background(), ScheduledCampaignAd{
		CampaignAd: CampaignAd{
			Placeholder: "placholder",
			Ratio:       0.5,
			Id:          "id2",
			Probability: 1,
			Fallback:    true,
			Ad: Ad{
				Source:      "source",
				Image:       "image",
				Link:        "http://link.com",
				Description: "desc",
				Company:     "company",
			},
		},
		Start: time.Now().Add(time.Hour * -1),
		End:   time.Now().Add(time.Hour),
	})
	assert.Nil(t, err)
	err = addOrUpdateUserTags(context.Background(), "1", []string{"javascript"})
	assert.Nil(t, err)
	_, err = db.Exec("insert into ad_tags (ad_id, tag) values ('id', ?), ('id2', ?)", "javascript", "php")
	assert.Nil(t, err)

	var res []CampaignAd
	res, err = fetchCampaigns(context.Background(), time.Now(), "1")
	dup := camp
	dup.IsTagTargeted = true
	assert.Nil(t, err)
	assert.Equal(t, []CampaignAd{dup}, res)
}

func TestFetchCampaignsWithBoth(t *testing.T) {
	migrateDatabase()
	initializeDatabase()
	defer tearDatabase()
	defer dropDatabase()

	err := addCampaign(context.Background(), ScheduledCampaignAd{
		CampaignAd: camp,
		Start:      time.Now().Add(time.Hour * -1),
		End:        time.Now().Add(time.Hour),
	})
	assert.Nil(t, err)
	err = addCampaign(context.Background(), ScheduledCampaignAd{
		CampaignAd: CampaignAd{
			Placeholder: "placholder",
			Ratio:       0.5,
			Id:          "id2",
			Probability: 1,
			Fallback:    true,
			Ad: Ad{
				Source:      "source",
				Image:       "image",
				Link:        "http://link.com",
				Description: "desc",
				Company:     "company",
			},
		},
		Start: time.Now().Add(time.Hour * -1),
		End:   time.Now().Add(time.Hour),
	})
	assert.Nil(t, err)
	err = setOrUpdateExperienceLevel(context.Background(), "1", "MORE_THAN_4_YEARS")
	assert.Nil(t, err)
	_, err = db.Exec("insert into ad_experience_level (ad_id, experience_level) values ('id', ?), ('id2', ?)", "MORE_THAN_4_YEARS", "MORE_THAN_6_YEARS")
	assert.Nil(t, err)

	var res []CampaignAd
	res, err = fetchCampaigns(context.Background(), time.Now(), "1")
	dup := camp
	dup.IsExpTargeted = true
	assert.Nil(t, err)
	assert.Equal(t, []CampaignAd{dup}, res)
}

func TestFetchCampaignsWithExperience(t *testing.T) {
	migrateDatabase()
	initializeDatabase()
	defer tearDatabase()
	defer dropDatabase()

	err := addCampaign(context.Background(), ScheduledCampaignAd{
		CampaignAd: camp,
		Start:      time.Now().Add(time.Hour * -1),
		End:        time.Now().Add(time.Hour),
	})
	assert.Nil(t, err)
	err = addCampaign(context.Background(), ScheduledCampaignAd{
		CampaignAd: CampaignAd{
			Placeholder: "placholder",
			Ratio:       0.5,
			Id:          "id2",
			Probability: 1,
			Fallback:    true,
			Ad: Ad{
				Source:      "source",
				Image:       "image",
				Link:        "http://link.com",
				Description: "desc",
				Company:     "company",
			},
		},
		Start: time.Now().Add(time.Hour * -1),
		End:   time.Now().Add(time.Hour),
	})
	assert.Nil(t, err)
	err = addOrUpdateUserTags(context.Background(), "1", []string{"javascript"})
	assert.Nil(t, err)
	_, err = db.Exec("insert into ad_tags (ad_id, tag) values ('id', ?), ('id2', ?)", "javascript", "php")
	assert.Nil(t, err)
	err = setOrUpdateExperienceLevel(context.Background(), "1", "MORE_THAN_4_YEARS")
	assert.Nil(t, err)
	_, err = db.Exec("insert into ad_experience_level (ad_id, experience_level) values ('id', ?), ('id2', ?)", "MORE_THAN_4_YEARS", "MORE_THAN_6_YEARS")
	assert.Nil(t, err)

	var res []CampaignAd
	res, err = fetchCampaigns(context.Background(), time.Now(), "1")
	dup := camp
	dup.IsExpTargeted = true
	dup.IsTagTargeted = true
	assert.Nil(t, err)
	assert.Equal(t, []CampaignAd{dup}, res)
}

func TestFetchCampaignsWithExperienceButNotMatching(t *testing.T) {
	migrateDatabase()
	initializeDatabase()
	defer tearDatabase()
	defer dropDatabase()

	err := addCampaign(context.Background(), ScheduledCampaignAd{
		CampaignAd: camp,
		Start:      time.Now().Add(time.Hour * -1),
		End:        time.Now().Add(time.Hour),
	})
	assert.Nil(t, err)
	err = addCampaign(context.Background(), ScheduledCampaignAd{
		CampaignAd: CampaignAd{
			Placeholder: "placholder",
			Ratio:       0.5,
			Id:          "id2",
			Probability: 1,
			Fallback:    true,
			Ad: Ad{
				Source:      "source",
				Image:       "image",
				Link:        "http://link.com",
				Description: "desc",
				Company:     "company",
			},
		},
		Start: time.Now().Add(time.Hour * -1),
		End:   time.Now().Add(time.Hour),
	})
	assert.Nil(t, err)
	err = addOrUpdateUserTags(context.Background(), "1", []string{"javascript"})
	assert.Nil(t, err)
	_, err = db.Exec("insert into ad_tags (ad_id, tag) values ('id', ?), ('id2', ?)", "javascript", "php")
	assert.Nil(t, err)
	err = setOrUpdateExperienceLevel(context.Background(), "1", "MORE_THAN_6_YEARS")
	assert.Nil(t, err)
	_, err = db.Exec("insert into ad_experience_level (ad_id, experience_level) values ('id', ?)", "MORE_THAN_4_YEARS")
	assert.Nil(t, err)

	var res []CampaignAd
	res, err = fetchCampaigns(context.Background(), time.Now(), "1")
	assert.Nil(t, err)
	assert.Equal(t, []CampaignAd(nil), res)
}

func TestFetchCampaignsWithExperienceButMissing(t *testing.T) {
	migrateDatabase()
	initializeDatabase()
	defer tearDatabase()
	defer dropDatabase()

	err := addCampaign(context.Background(), ScheduledCampaignAd{
		CampaignAd: camp,
		Start:      time.Now().Add(time.Hour * -1),
		End:        time.Now().Add(time.Hour),
	})
	assert.Nil(t, err)
	err = addCampaign(context.Background(), ScheduledCampaignAd{
		CampaignAd: CampaignAd{
			Placeholder: "placholder",
			Ratio:       0.5,
			Id:          "id2",
			Probability: 1,
			Fallback:    true,
			Ad: Ad{
				Source:      "source",
				Image:       "image",
				Link:        "http://link.com",
				Description: "desc",
				Company:     "company",
			},
		},
		Start: time.Now().Add(time.Hour * -1),
		End:   time.Now().Add(time.Hour),
	})
	assert.Nil(t, err)
	err = addOrUpdateUserTags(context.Background(), "1", []string{"javascript"})
	assert.Nil(t, err)
	_, err = db.Exec("insert into ad_tags (ad_id, tag) values ('id', ?), ('id2', ?)", "javascript", "php")
	assert.Nil(t, err)
	err = setOrUpdateExperienceLevel(context.Background(), "1", "MORE_THAN_6_YEARS")
	assert.Nil(t, err)

	var res []CampaignAd
	res, err = fetchCampaigns(context.Background(), time.Now(), "1")
	dup := camp
	dup.IsExpTargeted = false
	dup.IsTagTargeted = true
	assert.Nil(t, err)
	assert.Equal(t, []CampaignAd{dup}, res)
}
