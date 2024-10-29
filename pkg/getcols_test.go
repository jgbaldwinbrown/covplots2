package covplots

// import (
// 	"io"
// 	"testing"
// 	"strings"
// )
// 
// var intxt = `one	two	three
// four	five	six	seven
// eight	nine
// `
// 
// var outtxt = `two	three
// five	six
// nine	
// `
// 
// var hictxt = `chrom	start	end	hit_type	alt_hit_type	hits	alt_hits	pair_prop	alt_prop	pair_totprop	pair_totgoodprop	pair_totcloseprop	winsize	winstep	pair_fpkm	alt_fpkm	pair_prop_fpkm	alt_prop_fpkm	name
// 2L	-900	100	paired	self	4	441	0.008988764	0.99101124	3.741814e-09	7.2368536e-09	7.2368536e-09	1000	100	0.00036922835	0.040707426	0.008988764	0.99101124	Hybrid
// 2L	-800	200	paired	self	10	914	0.010822511	0.98917749	9.3545351e-09	1.8092134e-08	1.8092134e-08	1000	100	0.00092307088	0.084368678	0.010822511	0.98917749	Hybrid
// 2L	-700	300	paired	self	21	1638	0.012658228	0.98734177	1.9644524e-08	3.7993481e-08	3.7993481e-08	1000	100	0.0019384488	0.15119901	0.012658228	0.98734177	Hybrid
// 2L	-600	400	paired	self	33	2536	0.012845465	0.98715453	3.0869966e-08	5.9704042e-08	5.9704042e-08	1000	100	0.0030461339	0.23409078	0.012845465	0.98715453	Hybrid
// 2L	-500	500	paired	self	41	3463	0.011700913	0.98829909	3.8353594e-08	7.4177749e-08	7.4177749e-08	1000	100	0.0037845906	0.31965945	0.011700913	0.98829909	Hybrid
// 2L	-400	600	paired	self	43	3902	0.010899873	0.98910013	4.0224501e-08	7.7796176e-08	7.7796176e-08	1000	100	0.0039692048	0.36018226	0.010899873	0.98910013	Hybrid
// 2L	-300	700	paired	self	55	4821	0.011279737	0.98872026	5.1449943e-08	9.9506737e-08	9.9506737e-08	1000	100	0.0050768898	0.44501247	0.011279737	0.98872026	Hybrid
// `
// 
// var hicout = `chrom	start	end	hits
// 2L	-900	100	4
// 2L	-800	200	10
// 2L	-700	300	21
// 2L	-600	400	33
// 2L	-500	500	41
// 2L	-400	600	43
// 2L	-300	700	55
// `
// 
// func TestGetCols(t *testing.T) {
// 	r := strings.NewReader(intxt)
// 	r2 := GetCols(r, []int{1,2})
// 	var b strings.Builder
// 	io.Copy(&b, r2)
// 	if b.String() != outtxt {
// 		t.Errorf("b.String() %v != outtxt %v", b.String(), outtxt)
// 	}
// }
// 
// func TestGetMultipleCols(t *testing.T) {
// 	r1 := strings.NewReader(intxt)
// 	r2 := strings.NewReader(intxt)
// 
// 	frs := GetMultipleCols([]io.Reader{r1, r2}, []int{1,2})
// 
// 	var b1 strings.Builder
// 	var b2 strings.Builder
// 	done := make(chan struct{}, 2)
// 	go func () {
// 		io.Copy(&b1, frs[0])
// 		done <- struct{}{}
// 	}()
// 	go func () {
// 		io.Copy(&b2, frs[1])
// 		done <- struct{}{}
// 	}()
// 	<-done
// 	<-done
// 
// 	if b1.String() != outtxt {
// 		t.Errorf("b1.String() %v != outtxt %v", b1.String(), outtxt)
// 	}
// 	if b2.String() != outtxt {
// 		t.Errorf("b2.String() %v != outtxt %v", b2.String(), outtxt)
// 	}
// }
// 
// func TestHicPairCols(t *testing.T) {
// 	r1 := strings.NewReader(hictxt)
// 	r2 := strings.NewReader(hictxt)
// 
// 	frs, err := HicPairColumns([]io.Reader{r1, r2}, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 
// 	var b1 strings.Builder
// 	var b2 strings.Builder
// 	done := make(chan struct{}, 2)
// 	go func () {
// 		io.Copy(&b1, frs[0])
// 		done <- struct{}{}
// 	}()
// 	go func () {
// 		io.Copy(&b2, frs[1])
// 		done <- struct{}{}
// 	}()
// 	<-done
// 	<-done
// 
// 	if b1.String() != hicout {
// 		t.Errorf("b1.String() %v != hicout %v", b1.String(), hicout)
// 	}
// 	if b2.String() != hicout {
// 		t.Errorf("b2.String() %v != hicout %v", b2.String(), hicout)
// 	}
// }
