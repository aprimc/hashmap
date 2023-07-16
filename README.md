# hashmap

A generic hashmap and hashset for Go.

This is a weekend project. I'm not using it for anything important yet.

The goal of the project is to explore and find a usable interface for generic maps and sets in Go.

Goals:
1. Allow using slices as keys of a map.
2. Allow using sets as a keys of a map.
3. Allow using any type as a key of a map.

Choices:
1. Equality is defined on element type and not on the container.
2. Zero values are valid sets and maps.
3. Sets and maps behave as pointer types.
