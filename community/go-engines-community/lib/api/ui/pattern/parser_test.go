package pattern

import (
	"regexp/syntax"
	"slices"
	"testing"
)

func TestParseLiterals(t *testing.T) {
	testCases := []struct {
		pattern       string
		expected      [][]string
		desc          string
		expectedError bool
	}{
		{
			pattern:  "LITERAL",
			expected: [][]string{{"LITERAL"}},
			desc:     "Given standalone literal - should return single group with that literal",
		},
		{
			pattern:  "(LITERAL)",
			expected: [][]string{{"LITERAL"}},
			desc:     "Given literal in group - should return single group with that literal",
		},
		{
			pattern:  ".*",
			expected: [][]string{},
			desc:     "Given wildcard .* - should return no literals",
		},
		{
			pattern:  "(.*)",
			expected: [][]string{},
			desc:     "Given wildcard .* in group - should return no literals",
		},
		{
			pattern:  ".+",
			expected: [][]string{},
			desc:     "Given wildcard .+ - should return no literals",
		},
		{
			pattern:  "(.+)",
			expected: [][]string{},
			desc:     "Given wildcard .+ in group - should return no literals",
		},
		{
			pattern:       ".",
			expectedError: true,
			desc:          "Given single dot pattern - should return error",
		},
		{
			pattern:  "LITERAL.*",
			expected: [][]string{{"LITERAL"}},
			desc:     "Given literal with trailing .* - should return the literal",
		},
		{
			pattern:  "LITERAL.+",
			expected: [][]string{{"LITERAL"}},
			desc:     "Given literal with trailing .+ - should return the literal",
		},
		{
			pattern:  ".*LITERAL",
			expected: [][]string{{"LITERAL"}},
			desc:     "Given literal with leading .* - should return the literal",
		},
		{
			pattern:  ".+LITERAL",
			expected: [][]string{{"LITERAL"}},
			desc:     "Given literal with leading .+ - should return the literal",
		},
		{
			pattern:  "LITERAL(.*)",
			expected: [][]string{{"LITERAL"}},
			desc:     "Given literal with trailing (.*) group - should return the literal",
		},
		{
			pattern:  "LITERAL(.+)",
			expected: [][]string{{"LITERAL"}},
			desc:     "Given literal with trailing (.+) group - should return the literal",
		},
		{
			pattern:  "(.*)LITERAL",
			expected: [][]string{{"LITERAL"}},
			desc:     "Given literal with leading (.*) group - should return the literal",
		},
		{
			pattern:  "(.+)LITERAL",
			expected: [][]string{{"LITERAL"}},
			desc:     "Given literal with leading (.+) group - should return the literal",
		},
		{
			pattern:       "LITERAL.",
			expectedError: true,
			desc:          "Given literal with trailing single dot - should return error",
		},
		{
			pattern:       ".LITERAL",
			expectedError: true,
			desc:          "Given literal with leading single dot - should return error",
		},
		{
			pattern:       "LEFT_LITERAL.RIGHT_LITERAL",
			expectedError: true,
			desc:          "Given two literals with single dot between - should return error",
		},
		{
			pattern:       "LITERAL(.)",
			expectedError: true,
			desc:          "Given literal with trailing (.) group - should return error",
		},
		{
			pattern:       "(.)LITERAL",
			expectedError: true,
			desc:          "Given literal with leading (.) group - should return error",
		},
		{
			pattern:       "LEFT_LITERAL(.)RIGHT_LITERAL",
			expectedError: true,
			desc:          "Given two literals with (.) group between - should return error",
		},
		{
			pattern:  "LEFT_LITERAL.*RIGHT_LITERAL",
			expected: [][]string{{"LEFT_LITERAL"}, {"RIGHT_LITERAL"}},
			desc:     "Given two literals with .* between - should return both literals in separate groups",
		},
		{
			pattern:  "LEFT_LITERAL.+RIGHT_LITERAL",
			expected: [][]string{{"LEFT_LITERAL"}, {"RIGHT_LITERAL"}},
			desc:     "Given two literals with .+ between - should return both literals in separate groups",
		},
		{
			pattern:  "LEFT_LITERAL(.*)RIGHT_LITERAL",
			expected: [][]string{{"LEFT_LITERAL"}, {"RIGHT_LITERAL"}},
			desc:     "Given two literals with (.*) group between - should return both literals in separate groups",
		},
		{
			pattern:  "LEFT_LITERAL(.+)RIGHT_LITERAL",
			expected: [][]string{{"LEFT_LITERAL"}, {"RIGHT_LITERAL"}},
			desc:     "Given two literals with (.+) group between - should return both literals in separate groups",
		},
		{
			pattern:  "LEFT.*MIDDLE.+RIGHT",
			expected: [][]string{{"LEFT"}, {"MIDDLE"}, {"RIGHT"}},
			desc:     "Given three literals with wildcards between - should return all three literals in separate groups",
		},
		{
			pattern:  "LEFT(.*)MIDDLE(.+)RIGHT",
			expected: [][]string{{"LEFT"}, {"MIDDLE"}, {"RIGHT"}},
			desc:     "Given three literals with wildcard groups between - should return all three literals in separate groups",
		},
		{
			pattern:  "LOREM|IPSUM",
			expected: [][]string{{"LOREM", "IPSUM"}},
			desc:     "Given alternation of two literals - should return both literals in same group",
		},
		{
			pattern:  "(LOREM|IPSUM)",
			expected: [][]string{{"LOREM", "IPSUM"}},
			desc:     "Given alternation in group - should return both literals in same group",
		},
		{
			pattern:  "LOR(EM|AM)",
			expected: [][]string{{"LOREM", "LORAM"}},
			desc:     "Given literal prefix with suffix alternation - should return expanded combinations",
		},
		{
			pattern:  "(EM|AM)REM",
			expected: [][]string{{"EMREM", "AMREM"}},
			desc:     "Given prefix alternation with literal suffix - should return expanded combinations",
		},
		{
			pattern:  "[LM]OREM",
			expected: [][]string{{"LOREM", "MOREM"}},
			desc:     "Given character class prefix with literal suffix - should return expanded combinations",
		},
		{
			pattern:  "([LM])OREM",
			expected: [][]string{{"LOREM", "MOREM"}},
			desc:     "Given character class in group with literal suffix - should return expanded combinations",
		},
		{
			pattern:  "L[AO]REM",
			expected: [][]string{{"LAREM", "LOREM"}},
			desc:     "Given literal with middle character class - should return expanded combinations",
		},
		{
			pattern:  "L([AO])REM",
			expected: [][]string{{"LAREM", "LOREM"}},
			desc:     "Given literal with middle character class in group - should return expanded combinations",
		},
		{
			pattern:  "LORE[MN]",
			expected: [][]string{{"LOREM", "LOREN"}},
			desc:     "Given literal prefix with character class suffix - should return expanded combinations",
		},
		{
			pattern:  "LORE([MN])",
			expected: [][]string{{"LOREM", "LOREN"}},
			desc:     "Given literal prefix with character class in group - should return expanded combinations",
		},
		{
			pattern:  "(LO|LA)REM.*IPSUM",
			expected: [][]string{{"LAREM", "LOREM"}, {"IPSUM"}},
			desc:     "Given prefix alternation with suffix and wildcard separator - should return expanded first group and second literal",
		},
		{
			pattern:  "LOREM.*IPSU[MN]",
			expected: [][]string{{"LOREM"}, {"IPSUM", "IPSUN"}},
			desc:     "Given literal and wildcard separator with suffix character class - should return first literal and expanded second group",
		},
		{
			pattern:  "TEST[0-9].*",
			expected: [][]string{{"TEST0", "TEST1", "TEST2", "TEST3", "TEST4", "TEST5", "TEST6", "TEST7", "TEST8", "TEST9"}},
			desc:     "Given literal with digit range character class - should return all expanded combinations",
		},
		{
			pattern:  "(A|B)(C|D)",
			expected: [][]string{{"AC", "AD", "BC", "BD"}},
			desc:     "Given two consecutive alternations - should return cartesian product of combinations",
		},
		{
			pattern:  "^START.*END$",
			expected: [][]string{{"START"}, {"END"}},
			desc:     "Given anchored pattern with wildcard - should return both literals in separate groups",
		},
		{
			pattern: "^LOREM_(SED_IMPSUM|SED_UT_DOLOR|SED_UT_SIT)",
			expected: [][]string{
				{"LOREM_SED_IMPSUM", "LOREM_SED_UT_DOLOR", "LOREM_SED_UT_SIT"},
			},
			desc: "Given anchored prefix with suffix alternation - should return expanded combinations",
		},
		{
			pattern: "^LOREM.*_IMPSUMDOLOR.*|^LOREM.*_IMPSUMSIT.*|^LOREM.*_IMPSUMCONSECTETUR.*",
			expected: [][]string{
				{"LOREM"},
				{"_IMPSUMDOLOR", "_IMPSUMSIT", "_IMPSUMCONSECTETUR"},
			},
			desc: "Given top-level alternation with common prefix and middle literals - should return common prefix and combined middle literals",
		},
		{
			pattern: "^LOREM.*(SED_IMPSUM|SED_UT_DOLOR).*|^LOREM.*(SED_UT_DOLOR|SED_UT_SIT).*|^LOREM.*(SED_IMPSUM|SED_UT_SIT).*",
			expected: [][]string{
				{"LOREM"},
				{"SED_IMPSUM", "SED_UT_DOLOR", "SED_UT_SIT"},
			},
			desc: "Given top-level alternation with nested alternations - should return common prefix and merged alternation literals",
		},
		{
			pattern:       "[a-z]LOREM",
			expectedError: true,
			desc:          "Given character class with lowercase range - should return error due to group size limit",
		},
		{
			pattern:  "LOREM[A-C][0-2]",
			expected: [][]string{{"LOREMA0", "LOREMA1", "LOREMA2", "LOREMB0", "LOREMB1", "LOREMB2", "LOREMC0", "LOREMC1", "LOREMC2"}},
			desc:     "Given two consecutive character classes with small ranges - should return all cartesian combinations",
		},
		{
			pattern:       "[0-9][0-9]",
			expectedError: true,
			desc:          "Given two consecutive digit ranges - should return error due to group size limit",
		},
		{
			pattern:       "^LOREM0.*IPSUM0$|^LOREM1_IPSUM1$|^LOREM2_IPSUM2$|^LOREM3_IPSUM3$|^LOREM4_IPSUM4$|^LOREM5_IPSUM5$|^LOREM6_IPSUM6$|^LOREM7_IPSUM7$|^LOREM8_IPSUM8$|^LOREM9_IPSUM9$|^LOREM10_IPSUM10$",
			expectedError: true,
			desc:          "Given anchored pattern with wildcard separator - should return error due to group size limit",
		},
		{
			pattern:       "LOREM0.*LOREM1.*LOREM2.*LOREM3.*LOREM4.*LOREM5.*LOREM6.*LOREM7.*LOREM8.*LOREM9.*LOREM10",
			expectedError: true,
			desc:          "Given anchored pattern with wildcard separator - should return error due to group size limit",
		},
		{
			pattern:  "\\.",
			expected: [][]string{{"."}},
			desc:     "Given escaped dot - should return dot as literal",
		},
		{
			pattern:  "LO\\*REM",
			expected: [][]string{{"LO*REM"}},
			desc:     "Given escaped asterisk in literal - should return asterisk as part of literal",
		},
		{
			pattern:  "^LOREM",
			expected: [][]string{{"LOREM"}},
			desc:     "Given start anchor with literal - should return literal",
		},
		{
			pattern:  "LOREM",
			expected: [][]string{{"LOREM"}},
			desc:     "Given literal with end anchor - should return literal",
		},
		{
			pattern:  "^LOREM$",
			expected: [][]string{{"LOREM"}},
			desc:     "Given fully anchored literal - should return literal",
		},
		{
			pattern:  "((LOREM))",
			expected: [][]string{{"LOREM"}},
			desc:     "Given double nested capture groups - should return literal",
		},
		{
			pattern:  "[AB][CD][EF]",
			expected: [][]string{{"ACE", "ACF", "ADE", "ADF", "BCE", "BCF", "BDE", "BDF"}},
			desc:     "Given three consecutive two-element character classes - should return all 8 combinations",
		},
		{
			pattern:       "",
			expectedError: true,
			desc:          "Given empty pattern - should return error",
		},
		{
			pattern:       "L?OREM",
			expectedError: true,
			desc:          "Given literal with question mark - should return error",
		},
		{
			pattern:       "\\bLOREM\\b",
			expectedError: true,
			desc:          "Given literal with word boundary - should return error",
		},
		{
			pattern:       "\\BLOREM\\B",
			expectedError: true,
			desc:          "Given literal with non-word boundary - should return error",
		},
		{
			pattern:       "LOREM{2,3}",
			expectedError: true,
			desc:          "Given literal with repetition - should return error",
		},
		{
			pattern:       "LOREM*",
			expectedError: true,
			desc:          "Given literal with star - should return error",
		},
		{
			pattern:       "LOREM+",
			expectedError: true,
			desc:          "Given literal with plus - should return error",
		},
		{
			pattern:       "LONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXT",
			expectedError: true,
			desc:          "Given literal with long text - should return error",
		},
		{
			pattern:       "LONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEXTLONGTEX[S|T]",
			expectedError: true,
			desc:          "Given literal with long text and alternation - should return error",
		},
		{
			pattern: "LOREM.*IPSUM|IPSUM.*LOREM",
			expected: [][]string{
				{"LOREM"},
				{"IPSUM"},
			},
			desc: "Given alternation with same literal in both groups - should return both literals in same group without duplication",
		},
		{
			pattern: "LOREM.*(IPSUM|DOLOR)|IPSUM.*LOREM",
			expected: [][]string{
				{"LOREM"},
				{"IPSUM", "DOLOR"},
			},
			desc: "Given alternation with same literal in both groups - should return both literals in same group without duplication",
		},
		{
			pattern: "LOREM.*(IPSUM|DOLOR)|(IPSUM|DOLOR).*LOREM",
			expected: [][]string{
				{"LOREM"},
				{"IPSUM", "DOLOR"},
			},
			desc: "Given alternation with same literal in both groups - should return both literals in same group without duplication",
		},
		{
			pattern: "LOREM.*(IPSUM|DOLOR)|(IPSUM|SIT).*LOREM",
			expected: [][]string{
				{"LOREM", "SIT"},
				{"IPSUM", "DOLOR"},
			},
			desc: "Given alternation with same literal in both groups - should return both literals in same group without duplication",
		},
		{
			pattern: "LOREM.*(IPSUM|DOLOR)|(IPSUM|SIT).*(LOREM|AMET)",
			expected: [][]string{
				{"LOREM", "SIT"},
				{"IPSUM", "DOLOR", "AMET"},
			},
			desc: "Given alternation with same literal in both groups - should return both literals in same group without duplication",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			parsedTree, err := syntax.Parse(tc.pattern, syntax.Perl)
			if err != nil {
				t.Fatalf("Failed to parse pattern %q: %v", tc.pattern, err)
			}

			literals, err := parseLiterals(parsedTree)
			if err != nil {
				if !tc.expectedError {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			if err == nil {
				if tc.expectedError {
					t.Fatal("expected error, but got none")
				}
			}

			if len(literals) != len(tc.expected) {
				t.Errorf("Pattern %q: expected %d literals, got %d. Expected: %v, Got: %v",
					tc.pattern, len(tc.expected), len(literals), tc.expected, literals)
				return
			}

			for i, exp := range tc.expected {
				if !slices.Equal(exp, literals[i]) {
					t.Errorf("Pattern %q: expected literal %v at index %d, got %v", tc.pattern, exp, i, literals)
				}
			}
		})
	}
}

func BenchmarkParseLiterals_10Groups_10Each(b *testing.B) {
	pattern :=
		"(G0_0|G0_1|G0_2|G0_3|G0_4|G0_5|G0_6|G0_7|G0_8|G0_9).*" +
			"(G1_0|G1_1|G1_2|G1_3|G1_4|G1_5|G1_6|G1_7|G1_8|G1_9).*" +
			"(G2_0|G2_1|G2_2|G2_3|G2_4|G2_5|G2_6|G2_7|G2_8|G2_9).*" +
			"(G3_0|G3_1|G3_2|G3_3|G3_4|G3_5|G3_6|G3_7|G3_8|G3_9).*" +
			"(G4_0|G4_1|G4_2|G4_3|G4_4|G4_5|G4_6|G4_7|G4_8|G4_9).*" +
			"(G5_0|G5_1|G5_2|G5_3|G5_4|G5_5|G5_6|G5_7|G5_8|G5_9).*" +
			"(G6_0|G6_1|G6_2|G6_3|G6_4|G6_5|G6_6|G6_7|G6_8|G6_9).*" +
			"(G7_0|G7_1|G7_2|G7_3|G7_4|G7_5|G7_6|G7_7|G7_8|G7_9).*" +
			"(G8_0|G8_1|G8_2|G8_3|G8_4|G8_5|G8_6|G8_7|G8_8|G8_9).*" +
			"(G9_0|G9_1|G9_2|G9_3|G9_4|G9_5|G9_6|G9_7|G9_8|G9_9)"

	tree, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		b.Fatalf("parse failed: %v", err)
	}

	b.ResetTimer()

	for b.Loop() {
		_, err = parseLiterals(tree)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
