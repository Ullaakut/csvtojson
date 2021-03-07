package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

type Entries struct {
	Routines []Routine `json:"routines"`
}

type Routine struct {
	Title   string   `json:"title"`
	Sources []string `json:"sources"`
	Flags   []string `json:"flags"`
	Steps   Steps    `json:"steps"`
}

type Steps struct {
	Morning []Product `json:"Morning"`
	Evening []Product `json:"Evening"`
}

type Product struct {
	Name        string   `json:"name"`
	Link        string   `json:"link"`
	Type        string   `json:"type"`
	Ingredients []string `json:"ingredients"`
	Flags       []string `json:"flags"`
}

type intermRoutine struct {
	Title   string
	Sources map[string]struct{}
	Flags   map[string]struct{}
	Steps   Steps
}

type CSVProduct struct {
	Person string `csv:"person,omitempty"`

	MorningProduct            string `csv:"morning_product,omitempty"`
	MorningProductIngredients string `csv:"morning_product_ingredients,omitempty"`
	MorningProductType        string `csv:"morning_product_type,omitempty"`
	MorningProductLink        string `csv:"morning_product_link,omitempty"`

	NightProduct            string `csv:"night_product,omitempty"`
	NightProductIngredients string `csv:"night_product_ingredients,omitempty"`
	NightProductType        string `csv:"night_product_type,omitempty"`
	NightProductLink        string `csv:"night_product_link,omitempty"`

	Alcohol   bool `csv:"alcohol,omitempty"`
	SLS       bool `csv:"sls,omitempty"`
	Perfume   bool `csv:"perfume,omitempty"`
	VitaminC  bool `csv:"vitamin_c,omitempty"`
	VitaminB  bool `csv:"vitamin_b,omitempty"`
	Exfoliant bool `csv:"exfoliant,omitempty"`

	LinkMorning string `csv:"link_morning,omitempty"`
	LinkNight   string `csv:"link_night,omitempty"`
}

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Println("usage: csvtojson [input.csv] [output.json]")
		os.Exit(1)
	}

	file, err := os.OpenFile(args[0], os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var products []*CSVProduct
	if err := gocsv.UnmarshalFile(file, &products); err != nil {
		panic(err)
	}

	interm := make(map[string]intermRoutine)
	for _, p := range products {
		routine := interm[p.Person]
		routine.Title = p.Person

		if routine.Sources == nil {
			routine.Sources = make(map[string]struct{})
		}

		if p.LinkMorning != "" {
			routine.Sources[p.LinkMorning] = struct{}{}
		}
		if p.LinkNight != "" {
			routine.Sources[p.LinkNight] = struct{}{}
		}

		if routine.Flags == nil {
			routine.Flags = make(map[string]struct{})
		}

		if p.Alcohol {
			routine.Flags["A"] = struct{}{}
		}
		if p.SLS {
			routine.Flags["S"] = struct{}{}
		}
		if p.VitaminB {
			routine.Flags["B"] = struct{}{}
		}
		if p.VitaminC {
			routine.Flags["C"] = struct{}{}
		}
		if p.Exfoliant {
			routine.Flags["E"] = struct{}{}
		}
		if p.Perfume {
			routine.Flags["P"] = struct{}{}
		}

		if p.MorningProduct != "" {
			mp := Product{
				Name:        p.MorningProduct,
				Link:        p.MorningProductLink,
				Type:        p.MorningProductType,
			}

			for _, ingredient := range strings.Split(p.MorningProductIngredients, ",") {
				mp.Ingredients = append(mp.Ingredients, strings.TrimSpace(ingredient))
			}

			if p.Alcohol {
				mp.Flags = append(mp.Flags, "A")
			}
			if p.SLS {
				mp.Flags = append(mp.Flags, "S")
			}
			if p.VitaminB {
				mp.Flags = append(mp.Flags, "B")
			}
			if p.VitaminC {
				mp.Flags = append(mp.Flags, "C")
			}
			if p.Exfoliant {
				mp.Flags = append(mp.Flags, "E")
			}
			if p.Perfume {
				mp.Flags = append(mp.Flags, "P")
			}

			routine.Steps.Morning = append(routine.Steps.Morning, mp)
		}

		if p.NightProduct != "" {
			np := Product{
				Name:        p.NightProduct,
				Link:        p.NightProductLink,
				Type:        p.NightProductType,
			}

			for _, ingredient := range strings.Split(p.NightProductIngredients, ",") {
				np.Ingredients = append(np.Ingredients, strings.TrimSpace(ingredient))
			}

			if p.Alcohol {
				np.Flags = append(np.Flags, "A")
			}
			if p.SLS {
				np.Flags = append(np.Flags, "S")
			}
			if p.VitaminB {
				np.Flags = append(np.Flags, "B")
			}
			if p.VitaminC {
				np.Flags = append(np.Flags, "C")
			}
			if p.Exfoliant {
				np.Flags = append(np.Flags, "E")
			}
			if p.Perfume {
				np.Flags = append(np.Flags, "P")
			}

			routine.Steps.Evening = append(routine.Steps.Evening, np)
		}

		interm[p.Person] = routine
	}

	var entries Entries
	for title, iRoutine := range interm {
		var routine Routine

		routine.Title = title
		routine.Steps = iRoutine.Steps
		for flag := range iRoutine.Flags {
			routine.Flags = append(routine.Flags, flag)
		}
		for source := range iRoutine.Sources {
			routine.Sources = append(routine.Sources, source)
		}

		entries.Routines = append(entries.Routines, routine)
	}

	b, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile(args[1], b, 0o655); err != nil {
		panic(err)
	}
}
