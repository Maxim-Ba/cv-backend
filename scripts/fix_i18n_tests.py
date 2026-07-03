#!/usr/bin/env python3
"""Fix repository and services test files for i18n types."""
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def fix_technology_test(path: Path) -> None:
    text = path.read_text(encoding="utf-8")
    text = re.sub(
        r"// newPgText создает pgtype\.Text со значением\n"
        r"func newPgText\(s string\) pgtype\.Text \{\n"
        r"\treturn pgtype\.Text\{String: s, Valid: true\}\n"
        r"\}\n\n",
        "",
        text,
    )
    text = text.replace("Description: newPgText(", "Description: newNullableLocalizedText(")
    text = text.replace("Description: pgtype.Text{Valid: false}", "i18n.NullableLocalizedText{}")
    if "pkg/i18n" not in text:
        text = text.replace(
            '"github.com/jackc/pgx/v5/pgtype"',
            '"github.com/Maxim-Ba/cv-backend/pkg/i18n"\n\t"github.com/jackc/pgx/v5/pgtype"',
        )
    path.write_text(text, encoding="utf-8")


def fix_education_test(path: Path, add_i18n_import: bool = True) -> None:
    text = path.read_text(encoding="utf-8")
    text = re.sub(
        r'pgtype\.Text\{String: "([^"]*)", Valid: true\}',
        r'newNullableLocalizedText("\1")',
        text,
    )
    text = text.replace("pgtype.Text{Valid: false}", "i18n.NullableLocalizedText{}")
    text = re.sub(r'Course:\s+"([^"]+)"', r'Course: newLocalizedText("\1")', text)
    text = re.sub(r'Organization:\s+"([^"]+)"', r'Organization: newLocalizedText("\1")', text)
    text = text.replace("got.Name.String", "got.Name.Get(i18n.LocaleRU)")
    if add_i18n_import and "pkg/i18n" not in text:
        text = text.replace(
            '"github.com/jackc/pgx/v5/pgtype"\n',
            '"github.com/Maxim-Ba/cv-backend/pkg/i18n"\n',
        )
        text = text.replace('"github.com/jackc/pgx/v5/pgtype"', "")
    path.write_text(text, encoding="utf-8")


def fix_tag_test(path: Path) -> None:
    text = path.read_text(encoding="utf-8")
    text = re.sub(r'Name:\s+"([^"]+)"', r'Name: newLocalizedText("\1")', text)
    path.write_text(text, encoding="utf-8")


def fix_workhist_test(path: Path) -> None:
    text = path.read_text(encoding="utf-8")
    text = re.sub(r'Name:\s+"([^"]+)"', r'Name: newLocalizedText("\1")', text)
    text = re.sub(r'About:\s+"([^"]+)"', r'About: newLocalizedText("\1")', text)

    def repl_list(field: str, content: str) -> str:
        def inner(match: re.Match[str]) -> str:
            body = match.group(1).strip()
            if not body:
                return f"{field}: newLocalizedStringList()"
            return f"{field}: newLocalizedStringList({body})"

        return re.sub(rf"{field}:\s+\[\]string\{{([^}}]*)\}}", inner, content)

    text = repl_list("WhatIDid", text)
    text = repl_list("Projects", text)

    text = re.sub(
        r'Description: pgtype\.Text\{String: "([^"]+)", Valid: true\}',
        r'Description: newNullableLocalizedText("\1")',
        text,
    )
    text = text.replace("assert.Len(t, created.WhatIDid,", "assert.Len(t, created.WhatIDid.Get(i18n.LocaleRU),")
    text = text.replace("assert.Len(t, updated.WhatIDid,", "assert.Len(t, updated.WhatIDid.Get(i18n.LocaleRU),")
    text = re.sub(
        r'assert\.Equal\(t, "([^"]+)", updated\.WhatIDid\[0\]\)',
        r'assert.Equal(t, "\1", updated.WhatIDid.Get(i18n.LocaleRU)[0])',
        text,
    )
    text = re.sub(
        r"assert\.Equal\(t, \[\]string\{([^}]+)\}, result\.Content\[0\]\.WhatIDid\)",
        r"assert.Equal(t, newLocalizedStringList(\1), result.Content[0].WhatIDid)",
        text,
    )
    text = re.sub(
        r"created\.WhatIDid = \[\]string\{([^}]+)\}",
        r"created.WhatIDid = newLocalizedStringList(\1)",
        text,
    )
    if "pkg/i18n" not in text:
        text = text.replace(
            'entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"',
            'entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"\n\t"github.com/Maxim-Ba/cv-backend/pkg/i18n"',
        )
    path.write_text(text, encoding="utf-8")


def fix_services_education_test(path: Path) -> None:
    fix_education_test(path, add_i18n_import=True)
    text = path.read_text(encoding="utf-8")
    text = text.replace(
        "if result.Course != tt.mockEdu.Course {",
        "if result.Course.Get(i18n.LocaleRU) != tt.mockEdu.Course.Get(i18n.LocaleRU) {",
    )
    text = text.replace(
        "t.Errorf(\"Ожидался Course = %s, получили %s\", tt.mockEdu.Course, result.Course)",
        't.Errorf("Ожидался Course = %s, получили %s", tt.mockEdu.Course.Get(i18n.LocaleRU), result.Course.Get(i18n.LocaleRU))',
    )
    path.write_text(text, encoding="utf-8")


def fix_services_technology_test(path: Path) -> None:
    text = path.read_text(encoding="utf-8")
    text = re.sub(
        r'Description: pgtype\.Text\{String: "([^"]+)", Valid: true\}',
        r'Description: newNullableLocalizedText("\1")',
        text,
    )
    text = re.sub(
        r'LogoUrl:\s+pgtype\.Text\{String: "([^"]+)", Valid: true\}',
        r'LogoUrl: newPgText("\1")',
        text,
    )
    if "pkg/i18n" not in text:
        text = text.replace(
            '"github.com/jackc/pgx/v5/pgtype"\n\n',
            '"github.com/Maxim-Ba/cv-backend/pkg/i18n"\n\t"github.com/jackc/pgx/v5/pgtype"\n\n',
        )
    path.write_text(text, encoding="utf-8")


def fix_services_workhist_test(path: Path) -> None:
    text = path.read_text(encoding="utf-8")
    text = re.sub(r'Name:\s+"([^"]+)"', r'Name: newLocalizedText("\1")', text)
    text = re.sub(r'About:\s+"([^"]+)"', r'About: newLocalizedText("\1")', text)

    def repl_list(field: str, content: str) -> str:
        def inner(match: re.Match[str]) -> str:
            body = match.group(1).strip()
            if not body:
                return f"{field}: newLocalizedStringList()"
            return f"{field}: newLocalizedStringList({body})"

        return re.sub(rf"{field}:\s+\[\]string\{{([^}}]*)\}}", inner, content)

    text = repl_list("WhatIDid", text)
    text = repl_list("Projects", text)
    text = re.sub(
        r'LogoUrl:\s+pgtype\.Text\{String: "([^"]+)", Valid: true\}',
        r'LogoUrl: newPgText("\1")',
        text,
    )
    text = text.replace(
        "if len(result.WhatIDid) != len(tt.mockWH.WhatIDid) {",
        "if len(result.WhatIDid.Get(i18n.LocaleRU)) != len(tt.mockWH.WhatIDid.Get(i18n.LocaleRU)) {",
    )
    text = text.replace(
        "len(tt.mockWH.WhatIDid), len(result.WhatIDid)",
        "len(tt.mockWH.WhatIDid.Get(i18n.LocaleRU)), len(result.WhatIDid.Get(i18n.LocaleRU))",
    )
    if "pkg/i18n" not in text:
        text = text.replace(
            '"github.com/jackc/pgx/v5/pgtype"\n\n',
            '"github.com/Maxim-Ba/cv-backend/pkg/i18n"\n\t"github.com/jackc/pgx/v5/pgtype"\n\n',
        )
    path.write_text(text, encoding="utf-8")


def fix_services_tag_test(path: Path) -> None:
    text = path.read_text(encoding="utf-8")
    text = re.sub(r'Name:\s+"([^"]+)"', r'Name: newLocalizedText("\1")', text)
    if "pkg/i18n" not in text:
        text = text.replace(
            'entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"\n)',
            'entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"\n\t"github.com/Maxim-Ba/cv-backend/pkg/i18n"\n)',
        )
    path.write_text(text, encoding="utf-8")


def main() -> None:
    repo = ROOT / "internal" / "repository"
    svc = ROOT / "internal" / "services"

    fix_technology_test(repo / "technology_test.go")
    fix_education_test(repo / "education_test.go")
    fix_tag_test(repo / "tag_test.go")
    fix_workhist_test(repo / "workhist_test.go")

    fix_services_education_test(svc / "education_test.go")
    fix_services_technology_test(svc / "technology_test.go")
    fix_services_workhist_test(svc / "workhist_test.go")
    fix_services_tag_test(svc / "tag_test.go")

    print("Fixed i18n tests")


if __name__ == "__main__":
    main()
