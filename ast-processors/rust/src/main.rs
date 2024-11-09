use std::{io::{self, Read, Write}, process::{Command, Stdio}, thread};

use eyre::{eyre, OptionExt, Result};
use syn::Item;
use quote::ToTokens;

fn format_code(unformatted: &[u8], mut formatted: &mut Vec<u8>) -> Result<()> {
    let mut cmd = Command::new("rustfmt");
    cmd.arg("--edition").arg("2021");
    cmd.stdin(Stdio::piped()).stdout(Stdio::piped());

    let mut child = cmd.spawn()?;
    let mut stdin = child.stdin.take().ok_or_eyre("could not take stdin")?;
    let mut stdout = child.stdout.take().ok_or_eyre("could not take stdout")?;

    let unformatted = unformatted.to_vec();
    let handle = thread::spawn(move || stdin.write_all(&unformatted));

    io::copy(&mut stdout, &mut formatted)?;
    child.wait()?;

    match handle.join() {
        Ok(result) => result?,
        Err(_) => return Err(eyre!("stdin handle errored"))
    }

    Ok(())
}

fn main() {
    let data = &mut String::new();
    io::stdin().read_to_string(data).expect("could not read from stdin");

    let file = syn::parse_file(data).expect("could not read syntax");

    let functions = file.items.iter().filter_map(|item| if let Item::Fn(func) = item {
        Some(func)
    } else {
        None
    });

    let mut stdout = io::stdout();
    for function in functions {
        let string = function.into_token_stream().to_string();

        let mut formatted = Vec::new();
        format_code(string.as_bytes(), &mut formatted).expect("could not format code");

        let len = formatted.len().to_string();
        stdout.write_all(&[len.as_bytes(), &[b'|'], &formatted].concat()).expect("could not write to stdout");
    }
}
