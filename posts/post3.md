----
title: Test Post 3
description: Testing all kinds of markdown stuff
descriptionImage: https://github.com/nexfortisme/asset-store/blob/main/PlaceholderIcon1.png?raw=true
tags: [foo, test, Space Test, A fourth longer tag for some reason]
---


# Markdown Renderer Test Suite

This document is designed to test **common and advanced Markdown features**.

---

## 1. Text Formatting

**Bold text**  
*Italic text*  
***Bold + Italic***  
~~Strikethrough~~  
`Inline code`

> Blockquote  
> With multiple lines  
>> Nested blockquote

---

## 2. Headings

# H1
## H2
### H3
#### H4
##### H5
###### H6

---

## 3. Lists

### Unordered
- Item one
- Item two
  - Nested item
    - Deeply nested item
- Item three

### Ordered
1. First
2. Second
   1. Nested first
   2. Nested second
3. Third

### Task List
- [x] Completed task
- [ ] Incomplete task
- [ ] Another task

---

## 4. Links & References

Inline link: [Vue 3](https://vuejs.org)

Reference link: [Markdown Guide][md-guide]

[md-guide]: https://www.markdownguide.org/

Autolink: https://github.com

---

## 5. Images

Inline image:

![Vue Logo](https://upload.wikimedia.org/wikipedia/commons/9/95/Vue.js_Logo_2.svg)

Image with title:

![Vue Logo](https://upload.wikimedia.org/wikipedia/commons/9/95/Vue.js_Logo_2.svg "Vue.js Logo")

---

## 6. Code Blocks

### JavaScript
```js
export default {
  setup() {
    const message = "Hello Vue 3"
    return { message }
  }
}
