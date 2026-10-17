# Documentação em português

Este diretório contém a documentação em português do Atron.

**Estado no M0: vazio.** Nada foi traduzido ainda, e este arquivo não finge o
contrário.

## Aviso de autoridade

A documentação em inglês é a **autoritativa**. Se houver divergência entre um
texto deste diretório e o correspondente em inglês, vale o texto em inglês. Isso
está registrado no [ADR-0002](../decisions/0002-language-of-record.md).

O código, os comentários, as mensagens de diagnóstico e o texto principal da
documentação são em inglês. O guia mestre de implementação está em português por
ser o idioma do mantenedor, mas não substitui `docs/spec/`.

## O que será traduzido (e quando)

| Documento | Prioridade | Milestone |
|---|---|---|
| `README.md` | alta | M1 |
| Guia do usuário da linguagem | alta | junto do M8 |
| Referência da biblioteca padrão | média | M7 |
| Mensagens de diagnóstico (`ATRON_LANG=pt-BR`) | média | M6 |
| ADRs | baixa | não traduzir; são registro histórico |

## Como contribuir com traduções

1. Traduza a partir do texto em inglês vigente, nunca de uma versão antiga.
2. Mantenha os identificadores de código, nomes de comandos e códigos de
   diagnóstico (`E0101`, `R0012`) sem tradução.
3. Preserve a estrutura de títulos, para que os links entre documentos
   continuem funcionando.
4. Se o texto em inglês mudar, a tradução correspondente precisa ser revisada no
   mesmo pull request ou marcada como desatualizada no topo do arquivo.
