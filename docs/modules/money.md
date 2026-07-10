# Value Objects

## Money

The application represents monetary values using a dedicated `Money` value object instead of primitive numeric types.

### Representation

`Money` stores values using the smallest unit of the currency.

For INR:

- ₹1.00 → 100 paise
- ₹25.50 → 2550 paise
- ₹999.99 → 99999 paise

Internally, this is represented as an `int64`.

### Why not `float64`?

Floating-point numbers cannot precisely represent decimal values, which can lead to rounding errors.

Example:

```text
0.1 + 0.2 = 0.30000000000000004
```

Using integer arithmetic guarantees exact calculations for monetary values.

### Responsibilities

The `Money` value object is responsible for:

- Representing monetary values.
- Preventing invalid monetary values from being created.
- Encapsulating money-related behavior.

### Design

`Money` is implemented as a Value Object.

- It has no identity.
- Two `Money` values are equal if their amounts are equal.
- The internal amount is private and can only be accessed through its public API.

### Reusability

`Money` is a shared domain concept and is intended to be reused across multiple modules, including:

- Products
- Orders
- Discounts
- Payments

Future versions may extend `Money` with operations such as:

- Addition
- Subtraction
- Comparison
- String formatting
