# gomockctrl

This is a go linter based on `golang.org/x/tools/go/ast/inspector` that detects
a common `gomock` misconfiguration where `gomock.Controller` is intialized
outside `t.Run` scope where it's used. This leads to rare data races during
testing.  Example issue snippet:

``` go
    ctrl := gomock.NewController(t)
    fooService := NewMockFooService(ctrl)

    for _, tc := range testCases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()

            fooService.EXPECT().Foo()
```

The linter was constructed based on recommendations by ChatGPT.
