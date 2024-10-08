package celestial

import (
	"math"
	"testing"
	"time"
)

func nextEvent(t *testing.T, obs Observer, dt time.Time, event func(observer Observer, date time.Time) (time.Time, error)) time.Time {
	for offset := 0; offset < 366; offset++ {
		newdate := dt.Add(time.Duration(time.Duration(offset) * 24 * time.Hour))
		ti, err := event(obs, newdate)
		if err == nil {
			return ti
		}
	}

	t.Fatalf("Should not happen")
	return dt
}

func absDuration(n time.Duration) time.Duration {
	if n < 0 {
		return -n
	}
	return n
}

func diff(t1, t2 time.Time) time.Duration {
	return absDuration(t1.Sub(t2))
}

// TODO: this should probably have an own test
func almostEqualTime(t *testing.T, t1, t2 time.Time, allowedDiff time.Duration) {
	t.Helper()
	if d := diff(t1, t2); d > allowedDiff {
		t.Fatalf("diff: %q, t1 %q, t2 %q\n", d, t1, t2)
	}
}

func almostEqualFloat(t *testing.T, f1, f2, allowedDiff float64) {
	t.Helper()
	if abs := math.Abs(f1 - f2); abs > allowedDiff {
		t.Fatalf("diff: %f, f1 %f, f2 %f\n", abs, f1, f2)
	}
}

func TestNorwaySunUp(t *testing.T) {
	// """Test location in Norway where the sun doesn't set in summer."""
	june := time.Date(2019, 6, 5, 0, 0, 0, 0, time.UTC)
	obs := Observer{Latitude: 69.6, Longitude: 18.8, Elevation: 0.0}

	_, err := Sunrise(obs, june)
	if err == nil {
		t.FailNow()
	}
	_, err = Sunset(obs, june)
	if err == nil {
		t.FailNow()
	}

	// Find the next sunset and sunrise:
	nextSunrise := nextEvent(t, obs, june, Sunrise)
	nextSunset := nextEvent(t, obs, june, Sunset)

	if !nextSunrise.After(nextSunset) {
		t.FailNow()
	}
}

var london = Observer{Latitude: 51.509865, Longitude: -0.118092}

func TestDawn(t *testing.T) {
	type args struct {
		observer   Observer
		date       time.Time
		depression float64
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// Civil
		{args: args{observer: london, date: time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC), depression: DepressionCivil}, want: time.Date(2015, 12, 1, 7, 4, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 2, 0, 0, 0, 0, time.UTC), depression: DepressionCivil}, want: time.Date(2015, 12, 2, 7, 5, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 3, 0, 0, 0, 0, time.UTC), depression: DepressionCivil}, want: time.Date(2015, 12, 3, 7, 6, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 12, 0, 0, 0, 0, time.UTC), depression: DepressionCivil}, want: time.Date(2015, 12, 12, 7, 16, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 25, 0, 0, 0, 0, time.UTC), depression: DepressionCivil}, want: time.Date(2015, 12, 25, 7, 25, 0, 0, time.UTC)},
		// Nautical
		{args: args{observer: london, date: time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC), depression: DepressionNautical}, want: time.Date(2015, 12, 1, 6, 22, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 2, 0, 0, 0, 0, time.UTC), depression: DepressionNautical}, want: time.Date(2015, 12, 2, 6, 23, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 3, 0, 0, 0, 0, time.UTC), depression: DepressionNautical}, want: time.Date(2015, 12, 3, 6, 24, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 12, 0, 0, 0, 0, time.UTC), depression: DepressionNautical}, want: time.Date(2015, 12, 12, 6, 33, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 25, 0, 0, 0, 0, time.UTC), depression: DepressionNautical}, want: time.Date(2015, 12, 25, 6, 41, 0, 0, time.UTC)},
		// Astronomical
		{args: args{observer: london, date: time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC), depression: DepressionAstronomical}, want: time.Date(2015, 12, 1, 5, 41, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 2, 0, 0, 0, 0, time.UTC), depression: DepressionAstronomical}, want: time.Date(2015, 12, 2, 5, 42, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 3, 0, 0, 0, 0, time.UTC), depression: DepressionAstronomical}, want: time.Date(2015, 12, 3, 5, 44, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 12, 0, 0, 0, 0, time.UTC), depression: DepressionAstronomical}, want: time.Date(2015, 12, 12, 5, 52, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 25, 0, 0, 0, 0, time.UTC), depression: DepressionAstronomical}, want: time.Date(2015, 12, 25, 6, 1, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Dawn(tt.args.observer, tt.args.date, tt.args.depression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dawn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			almostEqualTime(t, got, tt.want, 60*time.Second)
		})
	}
}

func TestDusk(t *testing.T) {
	type args struct {
		observer   Observer
		date       time.Time
		depression float64
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// Civil
		{args: args{observer: london, date: time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC), depression: DepressionCivil}, want: time.Date(2015, 12, 1, 16, 34, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 2, 0, 0, 0, 0, time.UTC), depression: DepressionCivil}, want: time.Date(2015, 12, 2, 16, 34, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 3, 0, 0, 0, 0, time.UTC), depression: DepressionCivil}, want: time.Date(2015, 12, 3, 16, 33, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 12, 0, 0, 0, 0, time.UTC), depression: DepressionCivil}, want: time.Date(2015, 12, 12, 16, 31, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 25, 0, 0, 0, 0, time.UTC), depression: DepressionCivil}, want: time.Date(2015, 12, 25, 16, 36, 0, 0, time.UTC)},
		// Nautical
		{args: args{observer: london, date: time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC), depression: DepressionNautical}, want: time.Date(2015, 12, 1, 17, 16, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 2, 0, 0, 0, 0, time.UTC), depression: DepressionNautical}, want: time.Date(2015, 12, 2, 17, 16, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 3, 0, 0, 0, 0, time.UTC), depression: DepressionNautical}, want: time.Date(2015, 12, 3, 17, 16, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 12, 0, 0, 0, 0, time.UTC), depression: DepressionNautical}, want: time.Date(2015, 12, 12, 17, 14, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 25, 0, 0, 0, 0, time.UTC), depression: DepressionNautical}, want: time.Date(2015, 12, 25, 17, 19, 0, 0, time.UTC)},
		// Astronomical
		{args: args{observer: london, date: time.Date(2015, 12, 25, 0, 0, 0, 0, time.UTC), depression: DepressionAstronomical}, want: time.Date(2015, 12, 25, 17, 59, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2021, 30, 6, 0, 0, 0, 0, time.UTC), depression: DepressionAstronomical}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Dusk(tt.args.observer, tt.args.date, tt.args.depression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dawn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			almostEqualTime(t, got, tt.want, 60*time.Second)
		})
	}
}

func TestSunrise(t *testing.T) {
	type args struct {
		observer Observer
		date     time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{args: args{observer: london, date: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 1, 1, 8, 6, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 1, 7, 43, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 2, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 2, 7, 45, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 3, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 3, 7, 46, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 12, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 12, 7, 56, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 25, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 25, 8, 5, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sunrise(tt.args.observer, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dawn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			almostEqualTime(t, got, tt.want, 60*time.Second)
		})
	}
}

func TestSunset(t *testing.T) {
	type args struct {
		observer Observer
		date     time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{args: args{observer: london, date: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 1, 1, 16, 1, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 1, 15, 55, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 2, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 2, 15, 54, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 3, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 3, 15, 54, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 12, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 12, 15, 51, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 25, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 25, 15, 55, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sunset(tt.args.observer, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dawn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			almostEqualTime(t, got, tt.want, 60*time.Second)
		})
	}
}

func TestNoon(t *testing.T) {
	type args struct {
		observer Observer
		date     time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{args: args{observer: london, date: time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 1, 11, 49, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 2, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 2, 11, 49, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 3, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 3, 11, 50, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 12, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 12, 11, 54, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2015, 12, 25, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 25, 12, 00, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Noon(tt.args.observer, tt.args.date)
			almostEqualTime(t, got, tt.want, 60*time.Second)
		})
	}
}

func TestMidnight(t *testing.T) {
	type args struct {
		observer Observer
		date     time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{args: args{observer: london, date: time.Date(2016, 2, 18, 0, 0, 0, 0, time.UTC)}, want: time.Date(2016, 2, 18, 0, 14, 0, 0, time.UTC)},
		{args: args{observer: london, date: time.Date(2016, 10, 26, 0, 0, 0, 0, time.UTC)}, want: time.Date(2016, 10, 25, 23, 44, 0, 0, time.UTC)}, // TODO
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Midnight(tt.args.observer, tt.args.date)
			almostEqualTime(t, got, tt.want, 60*time.Second)
		})
	}
}

func TestTwilight(t *testing.T) {
	type args struct {
		observer  Observer
		date      time.Time
		direction SunDirection
	}
	tests := []struct {
		name      string
		args      args
		wantStart time.Time
		wantEnd   time.Time
		wantErr   bool
	}{
		// Rising
		{args: args{direction: SunDirectionRising, observer: london, date: time.Date(2019, 8, 29, 0, 0, 0, 0, time.UTC)}, wantStart: time.Date(2019, 8, 29, 4, 32, 0, 0, time.UTC), wantEnd: time.Date(2019, 8, 29, 5, 7, 0, 0, time.UTC)},
		// Setting
		{args: args{direction: SunDirectionSetting, observer: london, date: time.Date(2019, 8, 29, 0, 0, 0, 0, time.UTC)}, wantStart: time.Date(2019, 8, 29, 18, 54, 0, 0, time.UTC), wantEnd: time.Date(2019, 8, 29, 19, 30, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := Twilight(tt.args.observer, tt.args.date, tt.args.direction)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dawn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			almostEqualTime(t, start, tt.wantStart, 60*time.Second)
			almostEqualTime(t, end, tt.wantEnd, 60*time.Second)
		})
	}
}

var newDelhi = Observer{Latitude: 28.644800, Longitude: 77.216721}

func TestElevation(t *testing.T) {
	type args struct {
		observer   Observer
		date       time.Time
		refraction bool
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{args: args{refraction: true, observer: london, date: time.Date(2015, 12, 14, 11, 0, 0, 0, time.UTC)}, want: 14.381311},
		{args: args{refraction: true, observer: london, date: time.Date(2015, 12, 14, 20, 1, 0, 0, time.UTC)}, want: -37.3710156},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Elevation(tt.args.observer, tt.args.date, tt.args.refraction)
			// TODO: too far off. Python code uses accuracy of 0.001
			almostEqualFloat(t, got, tt.want, 0.005)
		})
	}
}

func TestAzimuth(t *testing.T) {
	type args struct {
		observer Observer
		date     time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{args: args{observer: london, date: time.Date(2015, 12, 14, 11, 0, 0, 0, time.UTC)}, want: 166.9676},
		{args: args{observer: london, date: time.Date(2015, 12, 14, 20, 1, 0, 0, time.UTC)}, want: 279.4093},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Azimuth(tt.args.observer, tt.args.date)
			// TODO: too far off. Python code uses accuracy of 0.001
			almostEqualFloat(t, got, tt.want, 0.01)
		})
	}
}

func TestZenith(t *testing.T) {
	type args struct {
		observer   Observer
		date       time.Time
		refraction bool
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		// {args: args{refraction: true, observer: london, date: time.Date(2019, 8, 29, 14, 34, 0, 0, time.UTC)}, want: 46}, // TODO: FIXME
		{args: args{refraction: true, observer: london, date: time.Date(2020, 2, 3, 10, 37, 0, 0, time.UTC)}, want: 71},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Zenith(tt.args.observer, tt.args.date, tt.args.refraction)
			almostEqualFloat(t, got, tt.want, 0.5)
		})
	}
}

func TestTimeAtElevation(t *testing.T) {
	type args struct {
		observer  Observer
		elevation float64
		date      time.Time
		direction SunDirection
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// Rising
		{args: args{direction: SunDirectionRising, elevation: 6, observer: london, date: time.Date(2016, 1, 4, 0, 0, 0, 0, time.UTC)}, want: time.Date(2016, 1, 4, 9, 5, 0, 0, time.UTC)},
		{args: args{direction: SunDirectionRising, elevation: 166, observer: london, date: time.Date(2016, 1, 4, 0, 0, 0, 0, time.UTC)}, want: time.Date(2016, 1, 4, 13, 20, 0, 0, time.UTC)},
		{args: args{direction: SunDirectionRising, elevation: 186, observer: london, date: time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC)}, want: time.Date(2015, 12, 1, 16, 34, 0, 0, time.UTC)},
		{args: args{direction: SunDirectionRising, elevation: -18, observer: london, date: time.Date(2016, 1, 4, 0, 0, 0, 0, time.UTC)}, want: time.Date(2016, 1, 4, 6, 0, 0, 0, time.UTC)},
		// Setting
		{args: args{direction: SunDirectionSetting, elevation: 14, observer: london, date: time.Date(2016, 1, 4, 0, 0, 0, 0, time.UTC)}, want: time.Date(2016, 1, 4, 13, 20, 0, 0, time.UTC)},
		// Error
		{wantErr: true, args: args{direction: SunDirectionRising, elevation: 20, observer: london, date: time.Date(2016, 1, 4, 0, 0, 0, 0, time.UTC)}, want: time.Date(2016, 1, 4, 6, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeAtElevation(tt.args.observer, tt.args.elevation, tt.args.date, tt.args.direction)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dawn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			almostEqualTime(t, got, tt.want, 180*time.Second)
		})
	}
}

func TestDaylight(t *testing.T) {
	type args struct {
		observer Observer
		date     time.Time
	}
	tests := []struct {
		name      string
		args      args
		wantStart time.Time
		wantEnd   time.Time
		wantErr   bool
	}{
		{args: args{observer: london, date: time.Date(2016, 1, 6, 0, 0, 0, 0, time.UTC)}, wantStart: time.Date(2016, 1, 6, 8, 5, 0, 0, time.UTC), wantEnd: time.Date(2016, 1, 6, 16, 7, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := Daylight(tt.args.observer, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dawn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			almostEqualTime(t, gotStart, tt.wantStart, 60*time.Second)
			almostEqualTime(t, gotEnd, tt.wantEnd, 60*time.Second)
		})
	}
}

func TestNight(t *testing.T) {
	type args struct {
		observer Observer
		date     time.Time
	}
	tests := []struct {
		name      string
		args      args
		wantStart time.Time
		wantEnd   time.Time
		wantErr   bool
	}{
		{args: args{observer: london, date: time.Date(2016, 1, 6, 0, 0, 0, 0, time.UTC)}, wantStart: time.Date(2016, 1, 6, 16, 46, 0, 0, time.UTC), wantEnd: time.Date(2016, 1, 7, 7, 25, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := Night(tt.args.observer, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dawn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			almostEqualTime(t, gotStart, tt.wantStart, 90*time.Second)
			almostEqualTime(t, gotEnd, tt.wantEnd, 60*time.Second)
		})
	}
}

func TestGoldenHour(t *testing.T) {
	type args struct {
		observer  Observer
		date      time.Time
		direction SunDirection
	}
	tests := []struct {
		name      string
		args      args
		wantStart time.Time
		wantEnd   time.Time
		wantErr   bool
	}{
		{args: args{direction: SunDirectionRising, observer: newDelhi, date: time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC)}, wantStart: time.Date(2015, 12, 1, 1, 10, 10, 0, time.UTC), wantEnd: time.Date(2015, 12, 1, 2, 0, 43, 0, time.UTC)},
		{args: args{direction: SunDirectionSetting, observer: london, date: time.Date(2016, 5, 18, 0, 0, 0, 0, time.UTC)}, wantStart: time.Date(2016, 5, 18, 19, 1, 0, 0, time.UTC), wantEnd: time.Date(2016, 5, 18, 20, 17, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := GoldenHour(tt.args.observer, tt.args.date, tt.args.direction)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dawn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			almostEqualTime(t, gotStart, tt.wantStart, 90*time.Second)
			almostEqualTime(t, gotEnd, tt.wantEnd, 60*time.Second)
		})
	}
}

func TestBlueHour(t *testing.T) {
	type args struct {
		observer  Observer
		date      time.Time
		direction SunDirection
	}
	tests := []struct {
		name      string
		args      args
		wantStart time.Time
		wantEnd   time.Time
		wantErr   bool
	}{
		{args: args{direction: SunDirectionRising, observer: london, date: time.Date(2016, 5, 19, 0, 0, 0, 0, time.UTC)}, wantStart: time.Date(2016, 5, 19, 3, 19, 0, 0, time.UTC), wantEnd: time.Date(2016, 5, 19, 3, 36, 0, 0, time.UTC)},
		{args: args{direction: SunDirectionSetting, observer: london, date: time.Date(2016, 5, 19, 0, 0, 0, 0, time.UTC)}, wantStart: time.Date(2016, 5, 19, 20, 18, 0, 0, time.UTC), wantEnd: time.Date(2016, 5, 19, 20, 35, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := BlueHour(tt.args.observer, tt.args.date, tt.args.direction)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dawn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			almostEqualTime(t, gotStart, tt.wantStart, 90*time.Second)
			almostEqualTime(t, gotEnd, tt.wantEnd, 60*time.Second)
		})
	}
}

func TestAdjustToHorizon(t *testing.T) {
	type args struct {
		elevation float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{args: args{elevation: 12000}, want: 3.517744168209966},
		{args: args{elevation: -1}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adjust_to_horizon(tt.args.elevation)
			almostEqualFloat(t, got, tt.want, 0.0000000000001)
		})
	}
}

func TestAdjustToObscuringFeature(t *testing.T) {
	type args struct {
		elevation0 float64
		elevation1 float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{args: args{elevation0: 0, elevation1: 100}, want: 0},
		{args: args{elevation0: 10, elevation1: 10}, want: 45},
		{args: args{elevation0: 3, elevation1: 4}, want: 53.130102354156},
		{args: args{elevation0: -10, elevation1: 10}, want: -45},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adjust_to_obscuring_feature(tt.args.elevation0, tt.args.elevation1)
			almostEqualFloat(t, got, tt.want, 0.0000000000001)
		})
	}
}
