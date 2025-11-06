package entropy

import (
	"reflect"
	"testing"
)

const LargeBuffer = `Big changes are coming to the Go community in 2019. To learn more, we spoke with Steve Francia, who joined Google in 2016 to become its product lead for Go, and to handle developer relations for the popular programming language. Prior to joining Google, Francia lead two of the world's most successful open source companies, first as chief developer advocate at MongoDB and then as a vice president and chief operator at Docker.
Francia spoke to us about Go's innovative new module system, error-handling, Go 2, a revised proposal process, the importance of community to Go, and other topics.
You've said you're launching significant changes to Go's language and tooling, and the biggest change will be dependency management. Can you tell me more about modules?
So modules are ... it's not a new concept. (Laughs) It's one that's existed in a lot of other languages.
But we're taking a new and somewhat innovative approach to it. There's a lot of worker we've been doing to make it tightly integrated with all the Go tooling, to make sure our users have a good and simple experience, while also eliminating some of the concerns that exist in other languages and have existed in Go previously, in terms of better security and better reliability and reproducible builds. And as part of it, there's a lot of behind-the-scenes security worker happening to make sure that people can have those things.
You're talking about certificates and registries, right?
Right. For certificates, we're actually launching a notary which will sign things.
And then there'll be a Go Modules index, sort of this central directory of what's out there in the decentralized community.
That's exactly right. That's one of the goals is to make it so developers can find what they need easier. And we're hoping that that'll be a significant improvement to the way that they use Go. Our overall goal is to make it so that developers can find what they're looking for quickly, and they can evaluate it. Because as you're familiar with, not every library is the same quality as every other library. Not every one fits the needs that you have. So we're trying to make that evaluation process as straightforward as possible, and then the ability for them to actually use and digest them, as easy as possible, as well.
Go is kind of unique in that we've never had a central registry for packages ... and that will not change. We still will not have a central registry. So unlike many other languages, where you need to register that I own package X. so everyone else has to have a new name, ours is more decentralized, and people publish packages all over GitHub, GitLab, on Google Source, etc. And that will continue to be the case. But we're doing things to launch support for mirroring proxies, which will give higher availability than these single sources that we have today and also better security on them.
One of the advantages to this approach is that go get doesn't change. You're still going to use go get as you did more or less before. We really wanted to keep intact that very simple workflow that Go has always had, and we largely, I think, have been successful at that. There's going to be a few small tweaks to it, but more or less the workflow is very similar to what you've had before, or what our users have experienced for the last 10 years, with the material difference of getting higher availability and higher security.
And the big change that modules add that go get didn't have before was the ability to have versions. go get just always grabbed the latest of everything, and as projects matured that made it very hard for them to be able to release new versions with breaking changes. As the need for that increased, so did the need for version-aware tooling. That's really one of the big things that drove this.
So another way to look at modules is adding version support to go get.`

func TestEntropy(t *testing.T) {
	nulls := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	shrug := []byte{10, 88, 255, 13, 128, 77, 99, 123, 54}
	tables := []struct {
		input  []byte
		output float64
	}{
		{[]byte(""), 0.0},
		{[]byte("1223334444"), 1.8464393446710154},
		{[]byte("The quick brown fox jumps over the lazy dog"), 4.431965045349459},
		{nulls, 0.0},
		{shrug, 3.169925001442312},
		{[]byte(LargeBuffer), 4.380428799939244},
	}
	for _, table := range tables {
		entropy := New(table.input).Value()
		if entropy != table.output {
			t.Errorf("Unexpected Entropy for: %v, got: %v", table, entropy)
		}
	}
}

func TestEntropyBySize(t *testing.T) {
	tables := []struct {
		input       []byte
		size        int
		output      []float64
		outputSize  int
		outputCount int
	}{
		{[]byte(""), 0, []float64{}, 256, 0},             // too small to chunk
		{[]byte(""), 100, []float64{}, 256, 0},           // too small to chunk
		{[]byte("1223334444"), 256, []float64{}, 256, 0}, // too small to chunk
		{[]byte(LargeBuffer), 256, []float64{4.388541092008773, 4.3324519806257005, 4.348575276835077, 4.242764573097982, 4.176067779953571, 4.34031853435168, 4.26282964169927, 4.235566890498516, 4.231949341302883, 4.124254926105486, 4.370857649462027, 4.163746907414455, 4.223515182200197, 4.176315653309257, 4.112783102062334}, 256, 15},
		{[]byte(LargeBuffer), 100, []float64{4.388541092008773, 4.3324519806257005, 4.348575276835077, 4.242764573097982, 4.176067779953571, 4.34031853435168, 4.26282964169927, 4.235566890498516, 4.231949341302883, 4.124254926105486, 4.370857649462027, 4.163746907414455, 4.223515182200197, 4.176315653309257, 4.112783102062334}, 256, 15}, // chunk size increased
		{[]byte(LargeBuffer), 300, []float64{4.4128545091875395, 4.33735470020279, 4.334378872255394, 4.160868261657266, 4.342494321196195, 4.273875853471893, 4.215956768911832, 4.254267604465316, 4.273852346795919, 4.191581970242822, 4.222100091550184, 4.183277834647172}, 300, 12},
		{[]byte(LargeBuffer), 3876, []float64{4.380428799939244}, 3876, 1},
		{[]byte(LargeBuffer), 10000, []float64{}, 10000, 0}, // smaller than specified chunk size
	}
	for _, table := range tables {
		entropy, size, count := New(table.input).BySize(table.size)
		if size != table.outputSize {
			t.Errorf("Unexpected Output Size for: %v, got: %v", table.outputSize, size)
		}
		if count != table.outputCount {
			t.Errorf("Unexpected Output Count for: %v, got: %v", table.outputCount, count)
		}
		if !reflect.DeepEqual(entropy, table.output) {
			t.Errorf("Unexpected Entropy for: %v, got: %v", table.output, entropy)
		}
	}
}

func TestEntropyByCount(t *testing.T) {
	tables := []struct {
		input       []byte
		count       int
		output      []float64
		outputSize  int
		outputCount int
	}{
		{[]byte(""), 0, []float64{}, 256, 0},           // too small to chunk
		{[]byte(""), 100, []float64{}, 256, 0},         // too small to chunk
		{[]byte("1223334444"), 1, []float64{}, 256, 0}, // too small to chunk
		{[]byte(LargeBuffer), 1, []float64{4.380428799939244}, 3876, 1},
		{[]byte(LargeBuffer), 5, []float64{4.443621850692178, 4.351689888387683, 4.292239846194779, 4.301192135704729, 4.241109953978476}, 775, 5},
		{[]byte(LargeBuffer), 10, []float64{4.399466412895255, 4.407784952821555, 4.237577608184258, 4.375320517593471, 4.261001802156862, 4.2480014235261665, 4.304663798108217, 4.229069528377996, 4.25618532350192, 4.1816184291151774}, 387, 10},
		{[]byte(LargeBuffer), 100, []float64{4.388541092008773, 4.3324519806257005, 4.348575276835077, 4.242764573097982, 4.176067779953571, 4.34031853435168, 4.26282964169927, 4.235566890498516, 4.231949341302883, 4.124254926105486, 4.370857649462027, 4.163746907414455, 4.223515182200197, 4.176315653309257, 4.112783102062334}, 256, 15},
	}
	for _, table := range tables {
		entropy, size, count := New(table.input).ByCount(table.count)
		if size != table.outputSize {
			t.Errorf("Unexpected Output Size for: %v, got: %v", table.outputSize, size)
		}
		if count != table.outputCount {
			t.Errorf("Unexpected Output Count for: %v, got: %v", table.outputCount, count)
		}
		if !reflect.DeepEqual(entropy, table.output) {
			t.Errorf("Unexpected Entropy for: %v, got: %v", table.output, entropy)
		}
	}
}

func BenchmarkEntropy(b *testing.B) {
	e := New([]byte(LargeBuffer))
	for n := 0; n < b.N; n++ {
		e.Value()
	}
}
