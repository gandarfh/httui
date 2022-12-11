<div align="center">
  <img
    src="assets/hero.png"
    alt="httui"
    style="width: 100%"
  />
  <p>httui is a Postman/Insomnia alternative</p>
</div>

## Introduction

httui was created to manage all my workspaces with rest apis and substitute applications like postman, insomnia, httpie.

## Example
https://user-images.githubusercontent.com/57275106/206926378-c9678fc3-d8e1-411e-9227-d2530de2b3ff.mp4

### Tech Stack

- [Go Lang](https://go.dev/)
- [Sqlite](https://www.sqlite.org/index.html)
- [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [Gorm](https://gorm.io/gorm)

### Features

- Multi workspaces

```
 ┌─────────┐ 1           N ┌────┐ 1         N ┌─────────┐
 │Workspace├──────────────►│Tags├────────────►│Resources│
 └─────────┘               └────┘             └─────────┘
```

- Fast rename (Rename only the name of workspace/tag/environment)
- Filter Resources
- Open in last resource opened
- Move resources between tags
- Custom environment to use in values
- Open response into nvim editor
- Update any information with nvim

## Contact

Twitter - [@gandarfh](https://twitter.com/gandarfh)
