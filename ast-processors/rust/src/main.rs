use std::io::{self, Read, Write};

use syn::{File, ImplItem, Item, ItemFn, ItemImpl};

fn impl_functions(impl_item: ItemImpl) -> Vec<Item> {
    impl_item.items.into_iter()
        .filter_map(|i| if let ImplItem::Fn(function) = i {
            Some(function)
        } else {
            None
        }).map(|i| ItemFn {
            attrs: i.attrs,
            sig: i.sig,
            vis: i.vis,
            block: Box::new(i.block)
        }).map(|i| i.into()).collect()
}

fn main() {
    let data = &mut String::new();
    io::stdin().read_to_string(data).expect("could not read from stdin");

    let file = syn::parse_file(data).expect("could not read syntax");
    let items = file.items.into_iter().filter_map(|item| Some(match item {
        Item::Fn(_) => vec![item],
        Item::Impl(impl_item) => impl_functions(impl_item),
        _ => return None
    })).flatten();

    let mut stdout = io::stdout();
    for item in items {
        let formatted = prettyplease::unparse(&File {
            shebang: None,
            attrs: vec![],
            items: vec![item]
        });

        let len = formatted.len().to_string();
        stdout.write_all(&[len.as_bytes(), &[b'|'], formatted.as_bytes()].concat()).expect("could not write to stdout");
    }
}
