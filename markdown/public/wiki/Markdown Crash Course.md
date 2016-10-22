# A Crash Course in Markdown

---

Well folks, you might have noticed that articles use a special type of scripting language for styling. This magical language is called **Markdown**! It involves converting patterns in *ordinary* text into 'styled' text. So here we go, I'll be laying down a few tips on how to express yourself using this language.

Before you begin, it would be good to let you know that the Markdown compiler that we use is called [blackfriday](https://github.com/russross/blackfriday). It is quite different from say for example GitHub's Markdown compiler. But nevertheless, with some tweaking and hacking you can perform all the functions of a markdown compiler and more!

## Headings

---

I am sure you have heard of headings in HTML (you guessed right, that's another kind of language. Markdown transpiles into HTML... If you are not sure what that means, it's absolutely OK! This deserves an article of its own :)).

It is completely alright if you are new to HTML or any of these weird languages that geeks like to invent, but in short HTML headings are: **H1**, **H2**, **H3**, **H4**, **H5** and **H6** denoted by `<Hx>` or `<hx>` where `x` can be anything from 1 to 6. Where **H1** is the largest in size and **H6** the smallest. I hope you have got the jist of it!

So, this is how you express headings in Markdown:

``` markdown
# Heading 1 or H1

## Heading 2

### H3

#### Grows smaller...

##### ... and even more smaller

###### And Voila! The tiniest size :3
```

Here is how that actually looks,

# Heading 1 or H1

## Heading 2

### H3

#### Grows smaller...

##### ... and even more smaller

###### And Voila! The tiniest size :3

Alright! So now you know how to use different headings, let us take a look at how to use a Horizontal rule.

## Horizontal rule

---

In HTML a Horizontal rule is created by the `<hr>` tag. In Markdown we use hyphens. You can use as many as you like, but 3 will suffice.

``` markdown
### This is a heading

---
```

Output:

### This is a heading

---

## Emphasis

---

Another interesting feature is the ability to make text bold, italic or even strikethrough text. To make text bold, you can put the text between `**` or `__` (two underscores). Let us check out an example:

``` markdown
**This makes text bold** Normal text

__but so does this!__
```

Output:

**This makes text bold** Normal text

__but so does this!__

Next up, *Italics*!

``` markdown
*Similar trick*

_except with just one astrick!_
```

Translates to:

*Similar trick*

_except with just one astrick!_

And finally, ~~strikethrough~~

``` markdown
~~cut in half~~
```

Output:

~~cut in half~~

Now let us look at lists.

## Lists

---

Lists can be both ordered and unordered.

An ordered list can be written as follows:

``` markdown
1. This is an ordered list
 1. Here is subitem 1.
 2. Notice the indentation. It can be anything, one space, two spaces or one tab, etc.
2. And another List
 1. Yet another subitem
```

Output:

1. This is an ordered list
 1. Here is subitem 1
 2. Notice the indentation. It can be anything, one space, two spaces or one tab, etc.
2. And another List
 1. Yet another subitem
 
Unordered lists can be defined as follows:
 
``` markdown
* Here is a bullet
 * indented one
* Yet another ordinary bullet
 * yet another indent
  * indent further...
* now, snap out of it! :D
```

Here is the output:

* Here is a bullet
 * indented one
* Yet another ordinary bullet
 * yet another indent
* and back...

It appears that on our Markdown compiler supports indentation upto one level only :(
Though, that is not something I would lose sleep over :D

## Links

---

Links are quite simple to use in Markdown even if they might look a bit intimidating! Using a link or an URL in HTML consists of two parts, the first being the text that a person sees when they see a link and of course the link's destination.

```
[Link #1](https://google.com)

[Link to Niec's Homepage](/)

[A link reference][ha! I am just a reference]

[References can be numbers as well!][1]

References are pretty cool! You can reuse them quite easily.

[ha! I am just a reference]: https://mozilla.org
[1]: http://reddit.com
```

Output:

[Link #1](https://google.com)

[Link to Niec's Homepage](/)

[A link reference][ha! I am just a reference]

[References can be numbers as well!][1]

References are pretty cool! You can reuse them.

[ha! I am just a reference]: https://mozilla.org
[1]: http://reddit.com

## Images

---

Here is a feature that most people would like to use in their articles: inserting images! To be honest, it is a lot like using Links, only slightly different. The place where the name of the link went, we use alternate text here in images (in case the image doesn't load).

```
Image without reference
![Kitty cat](/static/image/cat.jpeg)

With reference
![Adorable little chicken][ref]

[ref]: /static/image/chicks.jpg
```

Output looks like:

Image without reference
![Kitty cat](/static/image/cat.jpeg)

With reference
![Adorable little chicken][ref]

[ref]: /static/image/chicks.jpg