Label: About something
Description: Overwriting default description
+Keywords: about, more
My: custom param
IsCache: No
+++

# About

By default content files are cached.   
So you need to reload service
OR create `examples/.reload` file to trigger reload.

But this file uses param `IsCache: No`
and it's been reload every time content of file changes


Templates are cached forever. Only service stop/start will reload them.

## Image

![Logo](Mango.svg)

## URL

[Read more](/en/more)
