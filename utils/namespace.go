package utils

const XML_NS = "http://www.w3.org/xml/1998/namespace"

type namespace struct {
	Name  string
	Count int
}

type Namespaces struct {
	ns []namespace
}

func (n *Namespaces) findNS(name string) int {
	index := -1
	for i, el := range n.ns {
		if el.Name == name {
			return i
		}
		index = i
	}

	n.ns = append(n.ns, namespace{Name: name, Count: 0})
	return index + 1
}

func (n *Namespaces) Inc(name string) {
	index := n.findNS(name)
	n.ns[index].Count += 1
}

func (n *Namespaces) Dec(name string) {
	index := n.findNS(name)
	n.ns[index].Count -= 1
}

func (n *Namespaces) Has(name string) bool {
	index := n.findNS(name)
	return n.ns[index].Count > 0
}
