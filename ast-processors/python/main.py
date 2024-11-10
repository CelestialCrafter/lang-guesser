import ast
import sys

source = sys.stdin.read()
root = ast.parse(source)
for node in ast.iter_child_nodes(root):
    if not isinstance(node, ast.FunctionDef) and not isinstance(node, ast.AsyncFunctionDef):
        continue

    fn_source = ast.get_source_segment(source, node)
    if fn_source is None:
        continue

    sys.stdout.write(f'{len(fn_source) + 1}|')
    sys.stdout.write(fn_source)
