import ast
import sys

source = sys.stdin.read()
root = ast.parse(source)
for node in ast.iter_child_nodes(root):
    if not isinstance(node, ast.FunctionDef):
        continue

    fn_source = ast.get_source_segment(source, node)
    if fn_source is None:
        continue

    sys.stdout.write(f'{len(fn_source)}|')
    sys.stdout.write(fn_source)
