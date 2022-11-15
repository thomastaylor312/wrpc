use proc_macro::TokenStream;
use syn::{
    parse::{Parse, ParseStream, Result},
    punctuated::Punctuated,
    LitStr, Token,
};
use wit_bindgen_gen_guest_rust::Opts;

#[proc_macro]
pub fn generate(input: TokenStream) -> TokenStream {
    wit_bindgen_rust_macro_shared::generate::<Opt, Opts>(input, |opts| opts.build())
}

mod kw {
    syn::custom_keyword!(unchecked);
    syn::custom_keyword!(no_std);
    syn::custom_keyword!(raw_strings);
    syn::custom_keyword!(macro_export);
    syn::custom_keyword!(macro_call_prefix);
    syn::custom_keyword!(export_macro_name);
    syn::custom_keyword!(skip);
}

enum Opt {
    Unchecked,
    NoStd,
    RawStrings,
    MacroExport,
    MacroCallPrefix(LitStr),
    ExportMacroName(LitStr),
    Skip(Vec<LitStr>),
}

impl Parse for Opt {
    fn parse(input: ParseStream<'_>) -> Result<Self> {
        let l = input.lookahead1();
        if l.peek(kw::unchecked) {
            input.parse::<kw::unchecked>()?;
            Ok(Opt::Unchecked)
        } else if l.peek(kw::no_std) {
            input.parse::<kw::no_std>()?;
            Ok(Opt::NoStd)
        } else if l.peek(kw::raw_strings) {
            input.parse::<kw::raw_strings>()?;
            Ok(Opt::RawStrings)
        } else if l.peek(kw::macro_export) {
            input.parse::<kw::macro_export>()?;
            Ok(Opt::MacroExport)
        } else if l.peek(kw::macro_call_prefix) {
            input.parse::<kw::macro_call_prefix>()?;
            input.parse::<Token![:]>()?;
            Ok(Opt::MacroCallPrefix(input.parse()?))
        } else if l.peek(kw::export_macro_name) {
            input.parse::<kw::export_macro_name>()?;
            input.parse::<Token![:]>()?;
            Ok(Opt::ExportMacroName(input.parse()?))
        } else if l.peek(kw::skip) {
            input.parse::<kw::skip>()?;
            input.parse::<Token![:]>()?;
            let contents;
            syn::bracketed!(contents in input);
            let list = Punctuated::<_, Token![,]>::parse_terminated(&contents)?;
            Ok(Opt::Skip(list.iter().cloned().collect()))
        } else {
            Err(l.error())
        }
    }
}

impl wit_bindgen_rust_macro_shared::Configure<Opts> for Opt {
    fn configure(self, opts: &mut Opts) {
        match self {
            Opt::Unchecked => opts.unchecked = true,
            Opt::NoStd => opts.no_std = true,
            Opt::RawStrings => opts.raw_strings = true,
            Opt::MacroExport => opts.macro_export = true,
            Opt::MacroCallPrefix(prefix) => opts.macro_call_prefix = Some(prefix.value()),
            Opt::ExportMacroName(name) => opts.export_macro_name = Some(name.value()),
            Opt::Skip(list) => opts.skip.extend(list.iter().map(|i| i.value())),
        }
    }
}
