# bee-lite-java

[![Go Reference](https://pkg.go.dev/badge/github.com/Solar-Punk-Ltd/bee-lite-java.svg)](https://pkg.go.dev/github.com/Solar-Punk-Ltd/bee-lite-java)

This repository wraps [bee-lite](https://github.com/Solar-Punk-Ltd/bee-lite) functionality and
provides its core functionalities to implement Android applications what are able to use [Swarm network](https://docs.ethswarm.org/).
Due to **type restriction** this implemenentation also acts as an adapter to convert bridge between go and java by using the supported types.

## Development for Android platform

Please check the same section in this [README.md](https://github.com/Solar-Punk-Ltd/bee-lite/blob/master/README.md) it all applies here as well.

## Type restrictions

> Go-to-Android (gomobile) [Type Mapping Reference](https://pkg.go.dev/golang.org/x/mobile/cmd/gobind#hdr-Type_restrictions)

The following table outlines how Go types are mapped to Java/Kotlin when using `gomobile bind`.

| Go Type        | Java/Android Equivalent        | Category       |
| :------------- | :----------------------------- | :------------- |
| `bool`         | `boolean`                      | Basic Type     |
| `string`       | `String`                       | Basic Type     |
| `int`, `int32` | `int`                          | Numeric        |
| `int64`        | `long`                         | Numeric        |
| `float32`      | `float`                        | Numeric        |
| `float64`      | `double`                       | Numeric        |
| `uint8` (byte) | `byte`                         | Numeric        |
| `[]byte`       | `byte[]`                       | Buffer         |
| `error`        | `java.lang.Exception`          | Error Handling |
| `func`         | `interface` (Callback)         | Functional     |
| `struct`       | `class` (with getters/setters) | Object         |
| `interface`    | `interface`                    | Object         |

---

## ⚠️ Handling Unsupported Types

If your Go code uses types not listed above (like `uint64` or `[]string`), use the following workarounds:

### 1. Unsigned 64-bit Integers (`uint64`)

Java does not support unsigned 64-bit integers.

- **Go Side:** Cast to `int64` before returning.
- **Android Side:** Use `Long.toUnsignedString(value)` to display or `Long.compareUnsigned(...)` to compare.

### 2. String Slices (`[]string`)

`gomobile` does not support slices.

- **Workaround A (Wrapper):** Create a struct with `Get(i int) string` and `Length() int` methods.
- **Workaround B (JSON):** Serialize the slice to a JSON string in Go and parse it in Android.
- **Workaround C (Join):** Use `strings.Join(slice, "|")` in Go and `string.split("\\|")` in Kotlin/Java.

### 3. Maps and Complex Structs

For nested data or maps, it is recommended to use **JSON serialization** to pass data across the bridge.

## Compile with gomobile for Android

Target android api 21 is a nice sweet spot because it is Android Version: 5.0 (Lollipop) which is widely supported.

First you need gomobile and bind

```bash
make install
```

then

```bash
make build
```

OR if you want to do it your own

```bash
go install golang.org/x/mobile/cmd/gomobile@latest
go get golang.org/x/mobile/bind
```

In the root of this project run
`gomobile init`

Run the following:

`gomobile bind -target=android -androidapi=21 -o mobile.aar`
