# validator

[![Go Reference](https://pkg.go.dev/badge/github.com/siherrmann/validator.svg)](https://pkg.go.dev/github.com/siherrmann/validator)
[![Go Coverage](https://github.com/siherrmann/validator/wiki/coverage.svg)](https://raw.githack.com/wiki/siherrmann/validator/coverage.html)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/siherrmann/validator/blob/master/LICENSE)

Validator for structs and json maps to validate fields by tag.

# 💡 Goal of this package

The goal of the validator package is to provide a fast, simple validation tool that can validate existing structs and validate url values, json or yaml to update a struct with these values if they are valid.

Main focus is centralisation of validation to one object. You can obviously do one request object for each request type (which is sometimes still needed), but with this package you can do even something like that:

```go
type User struct {
    ID      int    `json:"id" del:"min1"`
    Name    string `json:"name" cre:"min1" upd:"min1, gr1min1"`
    Address string `json:"address" cre:"min1" upd:"min1, gr1min1"`
}
```

You can see the 3 different tags (`cre`, `upd`, `del`). With these tags you can use the same object in 3 different handlers with 3 different validations. In this example for creation you need name and address, for an update you need at least one of name or address, and for deletion you need an id greater 0.

---

# 🛠️ Installation

To integrate the validator package into your Go project, use the standard go get command:

```
go get github.com/siherrmann/validator
```

The validator needs Go 1.27 or newer, because it uses `encoding/json/v2`. In Go 1.27 the package is still behind the `jsonv2` experiment, which is enabled by default, so it does not build with `GOEXPERIMENT=nojsonv2`.

---

# 🚀 Getting started

To begin validating your structs, you either first create a new `Validator` instance and then invoke its `Validate` method with a pointer to your struct or use the wrapper function.

Here's a simple example demonstrating how to validate a User struct:

```go
type User struct {
    Name  string `vld:"min3"`
    Email string `vld:"rex^[^@]+@[^@]+\\.[^@]+$"`
    Age   int    `vld:"min18"`
}

user := User{
    Name:  "John Doe",
    Email: "john.doe@example.com",
    Age:   25,
}

// Example with validator
v := validator.NewValidator()
err := v.Validate(&user)
if err != nil {
    fmt.Println("Validation failed for user:", err)
}

// Example with wrapper functions (uses a new validator instance internally)
err = validator.Validate(&user)
if err != nil {
    fmt.Println("Validation failed for user:", err)
}
```

In these examples:

- `vld:"min3"`: Ensures the string field Name has a minimum length of 3 characters.
- `vld:"rex^[^@]+@[^@]+\\.[^@]+$"`: Validates Email against a regular expression pattern.
- `vld:"min18"`: Checks that the integer field Age has a minimum value of 18.

If you want to validate a request (json body or form values) and update for example the `User` struct only if the values are valid, you can use `UnmapOrUnmarshalAndValidate`.

```go
func CreateUser(w http.ResponseWriter, r *http.Request) {
    user := &User{}
    err := v.UnmapOrUnmarshalAndValidate(r, user)
    if err != nil {
        fmt.Println("Validation failed for new user:", err)
    }
    ...
}
```

If you want to use the `User` struct for multiple handlers like `CreateUser`, `UpdateUser` and `DeleteUser`, you could use custom tags like this:

```go
type User struct {
    ID    int    `delete:"min1" get:"min1"`
    Name  string `create:"min3" update:"min3, gr1min1"`
    Email string `create:"rex^[^@]+@[^@]+\\.[^@]+$" update:"rex^[^@]+@[^@]+\\.[^@]+$, gr1min1"`
    Age   int    `create:"min18" update:"min18, gr1min1"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    user := &User{}
    err := v.UnmapOrUnmarshalAndValidate(r, user, "create")
    if err != nil {
        fmt.Println("Validation failed for creating user:", err)
    }
    ...
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
    user := &User{}
    err := v.UnmapOrUnmarshalAndValidate(r, user, "update")
    if err != nil {
        fmt.Println("Validation failed for updating user:", err)
    }
    ...
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
    user := &User{}
    err := v.UnmapOrUnmarshalAndValidate(r, user, "delete")
    if err != nil {
        fmt.Println("Validation failed for deleting user:", err)
    }
    ...
}

func GetUser(w http.ResponseWriter, r *http.Request) {
    user := &User{}
    err := v.UnmapOrUnmarshalAndValidate(r, user, "get")
    if err != nil {
        fmt.Println("Validation failed for getting user:", err)
    }
    ...
}
```

You can also validate yaml (for example a config file) and update a struct with `UnmarshalYamlAndValidate`. To update a map by an array of `Validation` instead use `UnmarshalYamlAndValidateWithValidation`.

```go
type Config struct {
    Name string `yaml:"name" vld:"min3"`
    Port int    `yaml:"port" vld:"min1 && max65535"`
}

yamlInput, err := os.ReadFile("config.yaml")
if err != nil {
    ...
}

config := &Config{}
err = v.UnmarshalYamlAndValidate(yamlInput, config)
if err != nil {
    fmt.Println("Validation failed for config:", err)
}
```

---

# Validation

`Validate` validates a given struct by `vld` or custom tags. `ValidateAndUpdate` validates a json map and updates the given struct with the validated values. `UnmarshalAndValidate` and similar functions unpack the input (json body, form values or yaml), validate it and then update the given struct, so the struct is only changed if the input is valid. `ValidateWithValidation` and `ValidateAndUpdateWithValidation` do the same for a map by using an array of `Validation` (which is the equivalent for tags in a struct).

You can add a validate tag with the syntax `vld:"[requirement], [groups]"`.
Groups are separated by a space (eg. `gr1min1 gr2max1`).
Conditions and operators in a requirement are separated by a space (eg. `max0 || (min10 && max30)`).

All fields that you want to validate in the struct need a `vld` tag (or custom tag if specified).
If you don't want to validate the field you can add `vld:"-"`. If you then use an update function it does update it without validating.

The key of a field is taken from the `json` tag, falls back to the `yaml` tag and then to the field name. The option `omitempty` is read from the same tag. Fields tagged with `json:"-"` (or `yaml:"-"` without a `json` tag) are ignored, they are never validated or updated.

Embedded structs (also behind a pointer) without a name in the tag are flattened like in `encoding/json`, the same as fields with the option `inline` or `embed` (the name of the option in `encoding/json/v2`). Their fields are read, validated and updated on the level of the parent. If several fields have the same key, the shallowest one wins. At the same depth a field with a name in the tag wins, otherwise all of them are ignored. A nil embedded pointer is created on update if one of its fields is set, when validating a struct its fields are missing. Recursive types (eg. `Children []Node` in `Node`) are supported.

All functions use the same rules to read, validate and update a struct, so a struct behaves the same for json, yaml, form values and `Validate`. Values are converted the same way everywhere (eg. `"5"` into an `int` or a unix timestamp into a `time.Time`). Only fields whose type has an own `UnmarshalJSON` or whose `json` tag has the option `string` or `format:` are set with `encoding/json`. The json is decoded before that, so numbers above 2^53 have already lost precision.

JSON is parsed with `encoding/json/v2`, which is strict:

- Duplicate keys (eg. `{"name":"a","name":"b"}`) and invalid UTF-8 are rejected with an error instead of using the last value or replacing invalid characters. So the validated value is always the value every other JSON parser sees. This also applies to `istjson`.
- Keys are matched case-sensitive. `{"NAME":"a"}` does not set a field with the key `name`.
- A `null` sets the field to its zero value.

Types that can be unmarshaled from text (they implement `encoding.TextUnmarshaler`, e.g. `uuid.UUID` or `net.IP`) are validated as their string representation, also as pointers. So `istuuid` works on a `uuid.UUID` field and `min`/`max` check the length of the string (36 for a uuid). On update the string is parsed into the field type. Structs like `time.Time` are not affected by this. Keep in mind that the nil uuid (`00000000-0000-0000-0000-000000000000`) is a valid uuid, you can exclude it with `neq00000000-0000-0000-0000-000000000000`.

## Form values and slices

When url values (form values) are unpacked, a key with several values (eg. `tools=miner&tools=run`) becomes a list and a key with a single value (eg. `tools=miner`) becomes a lone value.
The struct field type decides the result: a `string` field gets the lone value, while a slice or array field (eg. `[]string`, `[]int`, `[2]string`) accepts either a list or a lone value, which is treated as a one-element list.
Validations on a slice field (eg. `min1`, `max3`, `con`, `frm`) also see a lone value as a list with one item.

```go
type Request struct {
    Tool  string   `json:"tool" vld:"min1"`
    Tools []string `json:"tools" vld:"min1"`
}
// tool=miner&tools=miner          -> Tool: "miner", Tools: []string{"miner"}
// tool=miner&tools=miner&tools=run -> Tool: "miner", Tools: []string{"miner", "run"}
```

## Requirement

You can build complex requirements by building a query of conditions, operators (`&&` (=AND), `||` (=OR) and `=>` (=IMPLIES)) and groups (with `(` and `)`).

A complex example for a password check (min length 8, max length 30, at least one capital letter, one small letter, one digit and one special character) would be:
`vld:"min8 max30 rex^(.*[A-Z])+(.*)$ rex^(.*[a-z])+(.*)$ rex^(.*\\d)+(.*)$ rex^(.*[\x60!@#$%^&*()_+={};':\"|\\,.<>/?~-])+(.*)$"`.
In this example all connections are `&&` (=AND) connections. Because behind the requirement check is a little parser you can also do more complex requirements with multiple conditions grouped and connected with AND and OR connections.
You can do for example `vld:"max0 || ((min10 && max30) || equTest)"` for a string that has to be either empty, the string `Test` or between 10 and 30 characters long. And yes, the outer brackets are not needed 😉.

`a => b` means "if `a` is fulfilled, `b` has to be fulfilled too". If `a` is not fulfilled, the requirement is fulfilled. `=>` binds weaker than `&&` and `||`, so `min1 && max5 => equ3 || equ4` is `(min1 && max5) => (equ3 || equ4)`. Use brackets to combine it with other conditions, eg. `(min3 => max5) && conb`. It is mostly used together with `rel` (see [Relations](#relations)).

### Condition types

Conditions have different usages per variable type:

- `-` - Not validating/update without validating.
- `equ` - `int/float/string == condition`, `len(array) == condition`
- `neq` - `int/float/string != condition`, `len(array) != condition`
- `min` - `int/float >= condition`, `len(strings.TrimSpace(string)/array/map) >= condition`, number of keys of a struct `>= condition`
- `max` - `int/float <= condition`, `len(strings.TrimSpace(string)/array/map) <= condition`, number of keys of a struct `<= condition`
- `con` - `strings.Contains(string, condition)`, `contains(array, condition)`, int/float ignored
- `nco` - `!strings.Contains(string, condition)`, `!contains(array, condition)`, int/float ignored
- `frm` - Checks if given comma separated list contains value/every item in array/every key in map.
- `nfr` - Checks if given comma separated list does not contain value/every item in array/every key in map.
- `rex` - `regexp.MatchString(condition, string)`, every item of an array/every key of a map (eg. env names). Numbers and bools are converted to a string first (`fmt.Sprintf("%d")` for int, `fmt.Sprintf("%v")` for float, eg. `3.5`).
- `ist` - Checks if string is of a certain format (e.g., `istuuid`, `istemail`, `istpassword`). Supported types: `uuid`, `unixmilli`, `timerfc3339`, `email`, `url`, `alpha`, `alphanum`, `numeric`, `hex`, `base64`, `json`, `jwt`, `ip`, `ipv4`, `ipv6`, `mac`, `password`.
- `def` - Default value, used if the key is missing or the value is invalid (eg. `equapple defunassigned`). On update the default is set in the struct.
- `cmn` - Checks if at least `condition` inner fields of a struct are correct (eg. `cmn2`), only usable on struct fields.
- `cmx` - Checks if at most `condition` inner fields of a struct are correct (eg. `cmx2`), only usable on struct fields.
- `unq` - Checks that every item of a list is unique (`unq`) or that a field of every item in a list of structs is unique (eg. `unqname`).
- `rel` - Checks the conditions in the brackets on another field (eg. `rel(delivery equshipping)`), see [Relations](#relations).
- `fun` - Checks the value with a custom function. The function has to be added to the validator (with `AddValidationFunc`, or with `AddValidationFuncContext` to read [injected values](#injected-values)), so it does not work with the wrapped functions. It can be used beside other requirements like `min3 && funYourCheckFunction`. This also allows you to check unsupported types by only using `funYourCheckFunction`.

A `null` value (e.g. `{"name": null}` in json or `name:` in yaml) counts as length/value `0` for `min` and `max`, so `max0` passes and `min1` fails. It is not equal to, does not contain, is not part of and does not match any condition, so `equ`, `con`, `frm`, `rex` and `ist` fail and `neq`, `nco` and `nfr` pass.

All conditions except `fun`, `def`, `ist`, `cmn`, `cmx` and `unq` can compare with the value of another field instead of a fixed value by using `$` and the path of the field (eg. `min$earliest_day`), see [Relations](#relations). With `${name}` they compare with a value injected for the validation call instead (eg. `equ${api_key}`), see [Injected values](#injected-values). A value starting with `$` is therefore always a reference, put it in quotes (eg. `equ'$100'`) to use it as a fixed value.

For `con` you need to put in a condition that is convertible to the underlying type of the array.
Eg. for an array of int the condition must be convertible to int (bad: `vld:"conA"`, good: `vld:"con1"`).

If you want to check more complex cases you can obviously replace `equ`, `neq`, `min`, `max` and `con` with one regular expression.

## Structs

The inner fields of a struct field are validated by their own tags. The requirement of the struct field itself works like this:

- Without a requirement (eg. `vld:"-"`) all inner fields have to be valid.
- With a requirement, it is checked on the struct itself, so `min` and `max` check the number of given keys. If the struct is empty (`{}` or `null`) and fulfills the requirement, the inner fields are not validated. If it is not empty, the inner fields have to be valid. A missing key is still an error, like for every other field.
- With `cmn` or `cmx` in the requirement, invalid inner fields are allowed. A field counts as correct if it is given (or has a default) and fulfills its own requirement. Invalid fields are left out, so on update they keep their current value.

```go
type Request struct {
    // Can be empty, but if not, the name has to be valid.
    Fruit struct {
        Name string `json:"name" vld:"min1"`
    } `json:"fruit" vld:"max0 || min1"`
    // Empty or at least two of the three fields have to be correct.
    Address struct {
        Street string `json:"street" vld:"min1"`
        City   string `json:"city" vld:"min1"`
        Zip    string `json:"zip" vld:"min5"`
    } `json:"address" vld:"max0 || cmn2"`
}
```

## Relations

Requirements can depend on other fields with `rel` and `$`.

- `rel(<path> <conditions>)` is fulfilled if the field at `<path>` fulfills `<conditions>`. Combined with `=>` this is a conditional requirement: `rel(delivery equshipping) => min1` means "if `delivery` is `shipping`, this field is required".
- `<condition>$<path>` compares with the value of the field at `<path>` instead of a fixed value, eg. `min$earliest_day`. With `frm` and `nfr` the referenced field can also be a list (eg. `frm$../items.name`) or a map, then its keys are used.

A path is relative to the level of the validated field:

- `name` - a field on the same level.
- `../name` - a field one level up (repeatable, eg. `../../name`). Items of a list of structs are one level below the struct with the list.
- `/name` - a field of the top level struct.
- `name.inner` - an inner field. On a list of structs it collects the field of every item (eg. `items.name`).

A missing field is `null` for `rel`, so `rel(other equyes)` is not fulfilled if `other` is missing. A missing field for `$` is an error.
A field with `rel` or `$` in its requirement is also validated if its key is missing or if it is empty with `omitempty`, the value is `null` (or `""`) then. So `rel(delivery equshipping) => min1` fails for a missing key if `delivery` is `shipping`. Missing fields are not added to the validated values.
References are resolved on the given input (not on validated or default values), so the order of the fields does not matter.
`rel` and `$` need the other fields, so they only work with `Validate`, `ValidateAndUpdate`, `ValidateWithValidation` and the functions using them, but not with `ValidateValueWithParser` or `RunValidatorsOnConditionGroup`.

```go
type Order struct {
    Delivery string `json:"delivery" vld:"frmpickup,shipping"`
    // Required if the order is shipped.
    Address string `json:"address,omitempty" vld:"rel(delivery equshipping) => min1"`
    // At least one item, every name only once.
    Items []struct {
        Name string `json:"name" vld:"min1"`
    } `json:"items" vld:"min1 && unqname"`
    // Optional, but the gift has to be one of the ordered items (one level up).
    Gift *struct {
        Item string `json:"item" vld:"frm$../items.name"`
    } `json:"gift,omitempty" vld:"-"`
    EarliestDay int `json:"earliest_day" vld:"min1"`
    // Not before the earliest day.
    LatestDay int `json:"latest_day" vld:"min$earliest_day"`
}
```

## Injected values

Values that are not part of the input (eg. the id or role of the logged in user, a limit of the account or a secret) can be injected for a single validation call with `validator.Inject`. They are stored in the context, so they only exist for the calls with this context (eg. one request) and concurrent requests do not see each others values.

```go
ctx := validator.Inject(r.Context(), "user_id", userID)
ctx = validator.Inject(ctx, "role", role)

// The unpack functions use the context of the request.
err := v.UnmarshalAndValidate(r.WithContext(ctx), &request)

// All other functions have a Context variant (eg. ValidateContext, ValidateAndUpdateContext, ValidateWithValidationContext).
err = v.ValidateContext(ctx, &request)
```

- `<condition>${name}` compares with the injected value (eg. `equ${user_id}`, `max${credit_limit}`). `frm` and `nfr` also accept a list or a map (eg. `frm${allowed_regions}`).
- `rel(${name} <conditions>)` checks the injected value (eg. `rel(${role} equadmin) => ...`). A missing value is `null`, like a missing field.
- Functions added with `AddValidationFuncContext` get the context and can read the values with `validator.Get`.

```go
v.AddValidationFuncContext(func(ctx context.Context, input any, astValue *model.AstValue) error {
    userID, ok := validator.Get[string](ctx, "user_id")
    if !ok {
        return errors.New("user id not injected")
    }
    // Check the input with the user id ...
    return nil
}, "ownedByUser")
```

A missing injected value for `<condition>${name}` is a configuration error (the value was not injected), so it is also returned inside an `||` group. Without a context (eg. `Validate` instead of `ValidateContext`) no values are injected.
For the security of injected values see [Injected values and secrets](#injected-values-and-secrets).

## Lists of structs

The inner fields of every item are validated by their own tags. The requirement of the list itself is checked on the list, eg. `min1` for at least one item or `unqname` for unique names.

## Groups

You also have the possibility to add groups. So if you want to check on an update, that at least one field is updated, you can add all fields to a group `upd:"min1, gr1min1"`.
Usable conditions for groups are `min` and `max`.

A small code example would be:

```go
type Error struct {
    ID                  int       `json:"id"`
    StatusCode          int       `json:"status_code" vld:"min100" upd:"min100, gr1min1"`
    Message             string    `json:"message" vld:"min1" upd:"min1, gr1min1"`
    UnderlyingException string    `json:"underlying_exception" vld:"min1, gr1min1" upd:"min1, gr1min1"`
    CreatedAt           time.Time `json:"created_at"`
}
```

`ID` and `CreatedAt` in this example are not getting validated, because you would probably not insert these but create these on database level.
`StatusCode` is necessary on creation and on update it is in a group with `Message` and `UnderlyingException` where one of them must be given.
One of `Message` and `UnderlyingException` is required on creation.

---

# 🔒 Security

The validator provides built-in security through tag-based field filtering. **Only fields with the specified validation tag will be updated**, protecting sensitive fields from unauthorized modification. This tag filtering also works recursively for nested structs and arrays of structs.

## Tag Behavior

This applies to all functions that update a struct (`ValidateAndUpdate`, `UnmarshalAndValidate` and the other unpack functions). With a custom tag (e.g., `"upd"`), fields are handled as follows:

```go
type User struct {
    ID        int    `json:"id"`                    // No tag: NEVER updated
    Password  string `json:"password"`              // No tag: NEVER updated
    Name      string `json:"name" upd:"-"`          // Ignore tag: Updated WITHOUT validation
    Email     string `json:"email" upd:"rex^[^@]+@[^@]+\\.[^@]+$"` // Full tag: Validated AND updated
    Role      string `json:"role"`                  // No tag: NEVER updated
}

// In your update handler:
user := &User{}
err := v.ValidateAndUpdate(jsonInput, user, "upd")
```

## Injected values and secrets

[Injected values](#injected-values) are often permissions or secrets, so they are handled with special care:

- **Not shared between calls**: The values belong to the context of the validation call, not to the validator. A validator used by concurrent requests never validates one request with the values of another one.
- **Never in error messages**: Validation errors are often returned to the client. An error of a condition with an injected value only names it (eg. `value does not fulfill equ${api_key}`), never the value itself. This also applies to `rel(${name} ...)`.
- **Constant time comparison**: `equ${name}` and `neq${name}` compare strings (and `[]byte` values) in constant time, so a secret can not be guessed by measuring how long the validation takes. Only the length of the secret can be measured.
- **Not cached**: The parsed requirements are cached per struct type, the injected values are only resolved during the call and never stored.

An example with permissions by role and a secret check:

```go
type UpdateDocument struct {
    // Users can only update their own documents, admins can update all.
    Owner string `json:"owner" vld:"rel(${role} equadmin) || equ${user_id}"`
    // Only admins can publish a document.
    Published bool `json:"published,omitempty" vld:"rel(${role} equadmin) || equfalse"`
    // The document has to be signed with the secret of the account.
    Signature string `json:"signature" vld:"funvalidSignature"`
}

v := validator.NewValidator()
v.AddValidationFuncContext(func(ctx context.Context, input any, astValue *model.AstValue) error {
    secret, ok := validator.Get[[]byte](ctx, "signing_secret")
    if !ok {
        return errors.New("signing secret not injected")
    }
    signature, _ := input.(string)
    if !verifySignature(signature, secret) { // eg. with crypto/hmac
        return errors.New("invalid signature")
    }
    return nil
}, "validSignature")

func handler(w http.ResponseWriter, r *http.Request) {
    ctx := validator.Inject(r.Context(), "user_id", currentUserID(r))
    ctx = validator.Inject(ctx, "role", currentRole(r))
    ctx = validator.Inject(ctx, "signing_secret", accountSecret(r))

    document := &UpdateDocument{}
    if err := v.UnmarshalAndValidate(r.WithContext(ctx), document); err != nil {
        // The error does not contain the injected user id, role or secret.
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
}
```

For a secret that is compared directly, use `equ${name}` instead of comparing it in a function, so the comparison is done in constant time.

---

# Testing

To run the tests run `go test ./...`.

## Benchmark

To run benchmarks run `go test -run '^$' -bench . -benchmem -count 10 > bench.txt`. To see the results in a nice way after the run install `go install golang.org/x/perf/cmd/benchstat@latest` and log the results to the console with `benchstat bench.txt`. To compare two versions, run the benchmarks for both into two files and use `benchstat old.txt new.txt`.

---

# ⭐ Features

- Validate json, yaml or form values by struct tag.
- Validate json, yaml or form values by validations.
- Validate and update a struct from `JsonMap`.
- **Tag-based validation**: Define validation rules directly within struct tags using a concise syntax (e.g. `vld:"(min3 && ncotest) || max0"`). Use multiple custom tags in one struct for multiple validation situations. The default tag is `vld`.
- **Relations**: Conditional requirements and comparisons with other fields, also on other levels (e.g. `vld:"rel(delivery equshipping) => min1"` or `vld:"min$earliest_day"`).
- **Injected values**: Validate against values of the current request like the user, the role or a secret (e.g. `vld:"rel(${role} equadmin) || equ${user_id}"`), safe for concurrent requests and without the values in error messages.
- **Custom validation functions**: Extend the validation capabilities by registering and using your own custom validation logic. Use it beside other requirements like `min3 && funYourCheckFunction`. This also allows you to check unsupported types by only using `funYourCheckFunction`.
- **Nested struct support**: Seamlessly validate complex data structures containing nested, embedded and recursive structs.
- **Array validation**: Apply validation rules to elements within arrays and slices.
- **Grouped validations**: Organize validation rules into logical groups for more granular control.
- **Advanced logical conditions**: Implement complex validation scenarios using logical operators (AND, OR and IMPLIES) within your tags.
