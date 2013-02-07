/*
Copyright 2011-2013 Paul Ruane.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package path

import (
	"path/filepath"
	"sort"
	"strings"
)

// Finds the root paths added to the tree
func Roots(paths []string) ([]string, error) {
	tree := buildTree(paths)

	roots := tree.roots()
	sort.Strings(roots)

	return roots, nil
}

// -

type tree struct {
	root *node
}

func (tree tree) roots() []string {
	roots := make([]string, 0, 100)
	return tree.root.findRoots(roots, "")
}

type node struct {
	name   string
	nodes  map[string]*node
	isRoot bool
}

func buildTree(paths []string) *tree {
	tree := tree{newNode("/", false)}

	for _, path := range paths {
		tree.add(path)
	}

	return &tree
}

func newNode(name string, root bool) *node {
	return &node{name, make(map[string]*node, 0), root}
}

func (tree *tree) add(path string) {
	pathParts := strings.Split(path, string(filepath.Separator))

	currentNode := tree.root
	partCount := len(pathParts)
	for index, pathPart := range pathParts {
		if pathPart == "" {
			pathPart = "/"
		}

		root := index == (partCount - 1)
		node, found := currentNode.nodes[pathPart]
		if !found {
			node = newNode(pathPart, root)
			currentNode.nodes[pathPart] = node
		} else {
			if !node.isRoot {
				if root {
					node.isRoot = true
				}
			}
		}

		currentNode = node
	}
}

func (node *node) findRoots(paths []string, path string) []string {
	path = filepath.Join(path, node.name)

	if node.isRoot {
		return append(paths, path)
	}

	for _, childNode := range node.nodes {
		paths = childNode.findRoots(paths, path)
	}

	return paths
}