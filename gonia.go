package gonia

import (
	"strconv"
	"os"
	"fmt"
	"bufio"
	"strings"
	"math"
	"sort"
)

type Config struct {
	Parsed	bool

	Score 	int64
	Mods	int

	Path 	string
}

type Gonia struct {
	PP		*PPS
	Stars		float64

	ConfParsed	 bool
	Conf 		*Config
	Map 		*Beatmap
}

type PPS struct {
	Total		float64
	Acc 		float64
	Strain 		float64
}

type Note struct {
	Exists		bool

	Key		int
	Start 		float64
	End 		float64
	Strain		float64
	HeldUntil	[]float64
	IndividualStrain []float64
}

type Beatmap struct {
	OD		float64
	Keys		int
	Notes		[]*Note
}

func ToInt16(s string) int16 {
	ret, err := strconv.ParseInt(s, 10, 16)

	if err != nil {
		return 0
	}

	return int16(ret)
}

func ToInt32(s string) int32 {
	ret, err := strconv.ParseInt(s, 10, 32)

	if err != nil {
		return 0
	}

	return int32(ret)
}

func ToInt64(s string) int64 {
	ret, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return 0
	}

	return ret
}

func ToFloat(s string) float32 {
	ret, err := strconv.ParseFloat(s, 32)

	if err != nil {
		return 0.0
	}

	return float32(ret)
}

func ToFloat64(s string) float64 {
	ret, err := strconv.ParseFloat(s, 32)

	if err != nil {
		return 0.0
	}

	return ret
}

func ToInt(s string) int {
	ret, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}

	return ret
}

func Trim(s string, trim string, pos string) string {
	if pos == "end" { 
		return strings.TrimSuffix(s, trim) 
	}

	return strings.TrimPrefix(s, trim) 
}

func (gonia *Gonia) ParseConf(args []string) (error) {
	if len(args) <= 0 {
		return fmt.Errorf("parse conf: no arguments to parse")
	}

	Conf := &Config{}

	// We will loop through all arguments
	// and see where the argument belong by
	// their suffixes and prefixes
	for i := 0; i < len(args); i++ {
		if strings.HasSuffix(args[i], "s") {
			Conf.Score = ToInt64(Trim(args[i], "s", "end"))
		} else if strings.HasSuffix(args[i], ".osu") {
			Conf.Path = args[i]
		} else if strings.HasPrefix(args[i], "+") {
			Conf.Mods = ToInt(Trim(args[i], "+", "start"))
		}
	}

	if Conf.Path == "" {
		return fmt.Errorf("parse error: no beatmap path was given.")
	}

	if Conf.Score > 1000000 {
		return fmt.Errorf("parse error: invalid score")
	}

	Conf.Parsed = true
	gonia.Conf = Conf

	return nil
}

func (gonia *Gonia) Parse(path string, score int64, mods int) (*Beatmap, error) {
	file, err := os.Open(path)
	
	if err != nil {
		fmt.Println("parse error: file not found")
		return nil, fmt.Errorf("parse error: %s", err)
	}

	if !gonia.ConfParsed {
		conf := &Config{}

		conf.Score = score
		conf.Mods = mods
		conf.Path = path

		gonia.Conf = conf
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	b := &Beatmap{}

	var Section string
	
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "[") {
			Section = line[1:len(line)-1]
			continue
		}

		switch Section {
		case "Difficulty":
			prop := strings.Split(line, ":")
			
			switch prop[0] {
			case "CircleSize":
				b.Keys = ToInt(prop[1])
			case "OverallDifficulty":
				b.OD = ToFloat64(prop[1])
			}
			break;
		case "HitObjects":
			n := &Note{}

			Note := strings.Split(line, ",")

			if len(Note) != 6 {
				fmt.Println("Invalid file format.")
				break
			}

			n.Key = int(math.Floor(ToFloat64(Note[0]) * float64(b.Keys) / 512))

			n.Start = ToFloat64(Note[2])
			n.End = ToFloat64(strings.Split(Note[5], ":")[0])
			n.Strain = 1 

			n.Exists = true

			if n.End == 0 {
				n.End = n.Start
			}

			for i := 0; i < b.Keys; i++ {
				n.IndividualStrain = append(n.IndividualStrain, 0)
				n.HeldUntil = append(n.HeldUntil, 0)
			}

			b.Notes = append(b.Notes, n)
			break;
		}
	}

	gonia.Map = b

	return b, nil
}

func (gonia *Gonia) CalculateStars() float64 {
	var TimeScale float64

	b := gonia.Map
	conf := gonia.Conf

	if (conf.Mods & (64 | 512)) > 0 {
		TimeScale = 1.5
	} else if (conf.Mods & 256) > 0 {
		TimeScale = 0.75
	} else {
		TimeScale = 1
	}
	
	var StrainStep float64 = 400.0 * TimeScale

	var WeightDecayBase float64 = 0.9
	var IndividualDecayBase float64 = 0.125
	var SrScalingFactor float64 = 0.018
	var OverallDecayBase float64 = 0.3

	var PrevNote *Note = b.Notes[0]

	for i := 0; i < len(b.Notes[1:]); i++ {
		Note := b.Notes[i]

		TimeElapsed := (Note.Start - PrevNote.Start) / TimeScale
		IndividualDecay := math.Pow(IndividualDecayBase, TimeElapsed / 1000.0)
		OverallDecay := math.Pow(OverallDecayBase, TimeElapsed / 1000.0)

		var HoldFactor float64 = 1.0
		var HoldAddition float64 = 0.0

		for j := 0; j < b.Keys; j++ {
			Note.HeldUntil[j] = PrevNote.HeldUntil[j]
			
			if Note.Start < Note.HeldUntil[j] && Note.End > Note.HeldUntil[j] {
				HoldAddition = 1.0
			} else if Note.End == Note.HeldUntil[j] {
				HoldAddition = 0.0
			} else if Note.End < Note.HeldUntil[j] {
				HoldFactor = 1.25
			}
			Note.IndividualStrain[j] = PrevNote.IndividualStrain[j] * IndividualDecay
		}

		Note.HeldUntil[Note.Key] = Note.End
		Note.IndividualStrain[Note.Key] += 2.0 * HoldFactor

		Note.Strain = PrevNote.Strain * OverallDecay + (1.0 + HoldAddition) * HoldFactor

		PrevNote = Note
	}

	var StrainTable []float64
	IntervalEndTime := StrainStep
	var MaximumStrain float64
	var PreviousNote Note

	for i := 0; i < len(b.Notes); i++ {
		Note := b.Notes[i]

		for Note.Start > IntervalEndTime {
			StrainTable = append(StrainTable, MaximumStrain)

			if !PreviousNote.Exists {
				IntervalEndTime += StrainStep
				continue
			}

			IndividualDecay := math.Pow(IndividualDecayBase, (IntervalEndTime - PreviousNote.Start) / 1000.0)
			OverallDecay := math.Pow(OverallDecayBase, (IntervalEndTime - PreviousNote.Start) / 1000.0)
			MaximumStrain = Note.IndividualStrain[Note.Key] * IndividualDecay + Note.Strain * OverallDecay

			IntervalEndTime += StrainStep
		}

		MaximumStrain = math.Max(Note.IndividualStrain[Note.Key] + PreviousNote.Strain, MaximumStrain)

		PreviousNote = *Note
	}

	var Difficulty float64
	var Weight float64 = 1.0
	
	sort.Slice(StrainTable, func(i, j int) bool {
		return StrainTable[i] > StrainTable[j]
	})	

	for i := 0; i < len(StrainTable); i++ {
		Difficulty += StrainTable[i] * Weight
		Weight *= WeightDecayBase
	}

	gonia.Stars = Difficulty * SrScalingFactor

	return gonia.Stars
}

func (gonia *Gonia) ComputeStrainValue() float64 {
	var StrainValue float64 = math.Pow(5 * math.Max(1, gonia.Stars / 0.2) - 4.0, 2.2) / 135.0

	StrainValue *= (1.0 + 0.1 * math.Min(1.0, float64(len(gonia.Map.Notes)) / 1500.0))

	score := gonia.Conf.Score

	if score <= 500000 {
		StrainValue = 0.0
	} else if score <= 600000 {
		StrainValue *= ((float64(score) - 500000.0) / 100000.0 * 0.3)
	} else if score <= 700000 {
		StrainValue *= (0.3 + (float64(score) - 600000.0) / 100000.0 * 0.25)
	} else if score <= 800000 {
		StrainValue *= (0.55 + (float64(score) - 700000.0) / 100000.0 * 0.20)
	} else if score <= 900000 {
		StrainValue *= (0.75 + (float64(score) - 800000.0) / 100000.0 * 0.15)
	} else {
		StrainValue *= (0.9 + (float64(score) - 900000) / 100000.0 * 0.1)
	}

	return StrainValue
}

func (gonia *Gonia) GetHitWindow300() int {
	var HitWindow300 float64 = 34 + 3 * (math.Min(10, math.Max(0, 10 - gonia.Map.OD)))

	if (gonia.Conf.Mods & 16) > 0 {
		HitWindow300 /= 1.4
	} else if (gonia.Conf.Mods & 2) > 0 {
		HitWindow300 *= 1.4
	}

	if (gonia.Conf.Mods & 64) > 0 {
		HitWindow300 *= 1.5
	} else if (gonia.Conf.Mods & 256) > 0 {
		HitWindow300 *= 0.75
	}

	return int(HitWindow300)
}

func (gonia *Gonia) ComputeAccValue() float64 {
	HitWindow300 := gonia.GetHitWindow300()

	if HitWindow300 <= 0 {
		return 0.0
	}

	return math.Max(0, 0.2 - (float64(HitWindow300 - 34) * 0.006667)) * gonia.PP.Strain * math.Pow(math.Max(0.0, float64(gonia.Conf.Score - 960000) / 40000.0), 1.1)
}

func (gonia *Gonia) CalculatePP() {
	var Multiplier float64 = 0.8
	var ScoreMultiplier float64 = 1.0

	if (gonia.Conf.Mods & 1) > 0 {
		Multiplier *= 0.9
		ScoreMultiplier *= 0.5
	}

	if (gonia.Conf.Mods & 4096) > 0 {
		Multiplier *= 0.95
	}

	if (gonia.Conf.Mods & 2) > 0 {
		Multiplier *= 0.5
		ScoreMultiplier *= 0.5
	}

	if (gonia.Conf.Mods & 256) > 0 {
		ScoreMultiplier *= 0.5
	}

	gonia.Conf.Score *= int64(1.0 / ScoreMultiplier)

	gonia.PP = &PPS{}
	gonia.PP.Strain = gonia.ComputeStrainValue()
	gonia.PP.Acc = gonia.ComputeAccValue()
	gonia.PP.Total = math.Pow(math.Pow(gonia.PP.Strain, 1.1) + math.Pow(gonia.PP.Acc, 1.1), 1.0 / 1.1) * Multiplier
}
