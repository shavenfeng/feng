package feng

import (
	"strings"
)

type trees map[string]*node

type node struct {
	pattern  string
	handlers []HandlerFunc
	children []*node
}

func (n *node) addNode(method, path string, handlers ...HandlerFunc) {
	var (
		childNode *node
	)
	patterns := strings.Split(strings.TrimPrefix(path, "/"), "/")
	pattern := patterns[0]
	if n.children == nil {
		childNode = &node{
			pattern: pattern,
		}
	} else {
		for _, child := range n.children {
			if pattern == child.pattern {
				child.addNode(method, strings.Join(patterns[1:], "/"), handlers...)
				return
			}
			childNode = &node{
				pattern: pattern,
			}
		}
	}
	n.children = append(n.children, childNode)
	if len(patterns) > 1 {
		childNode.addNode(method, strings.Join(patterns[1:], "/"), handlers...)
	} else {
		childNode.handlers = handlers
	}
}

func (n *node) findNode(method string, path string, params Params) *node {
	patterns := strings.Split(path, "/")
	firstPattern := patterns[0]
	if firstPattern == "" {
		patterns[0] = "/"
	}
	if len(patterns) == 1 && n.pattern == firstPattern {
		return n
	}
	if (n.pattern == firstPattern || strings.HasPrefix(n.pattern, ":")) && len(patterns) == 1 {
		return n
	}
	children := n.children
	for _, keyNode := range children {
		newPath := strings.Join(patterns[1:], "/")
		if patterns[1] == keyNode.pattern {
			return keyNode.findNode(method, newPath, params)
		}
		if keyNode != nil && strings.HasPrefix(keyNode.pattern, ":") {
			if params != nil {
				params[keyNode.pattern[1:]] = patterns[1]
			}
			return keyNode.findNode(method, newPath, params)
		}
	}
	return nil
}
