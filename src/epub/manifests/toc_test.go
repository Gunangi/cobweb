package manifests_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/cloakwiss/cobweb/epub/manifests"
)

func TestDirectoryTree(t *testing.T) {
	files := []string{
		"/META-INF/com.apple.ibooks.display-options.xml",
		"/META-INF/container.xml",
		"/OEBPS/Access.xhtml",
		"/OEBPS/Agent.xhtml",
		"/OEBPS/Application.xhtml",
		"/OEBPS/ArgumentError.xhtml",
		"/OEBPS/ArithmeticError.xhtml",
		"/OEBPS/Atom.xhtml",
		"/OEBPS/BadArityError.xhtml",
		"/OEBPS/BadBooleanError.xhtml",
		"/OEBPS/BadFunctionError.xhtml",
		"/OEBPS/BadMapError.xhtml",
		"/OEBPS/BadStructError.xhtml",
		"/OEBPS/Base.xhtml",
		"/OEBPS/Behaviour.xhtml",
		"/OEBPS/Bitwise.xhtml",
		"/OEBPS/Calendar.ISO.xhtml",
		"/OEBPS/Calendar.TimeZoneDatabase.xhtml",
		"/OEBPS/Calendar.UTCOnlyTimeZoneDatabase.xhtml",
		"/OEBPS/Calendar.xhtml",
		"/OEBPS/CaseClauseError.xhtml",
		"/OEBPS/Code.Fragment.xhtml",
		"/OEBPS/Code.LoadError.xhtml",
		"/OEBPS/Code.xhtml",
		"/OEBPS/Collectable.xhtml",
		"/OEBPS/CompileError.xhtml",
		"/OEBPS/CondClauseError.xhtml",
		"/OEBPS/Config.Provider.xhtml",
		"/OEBPS/Config.Reader.xhtml",
		"/OEBPS/Config.xhtml",
		"/OEBPS/Date.Range.xhtml",
		"/OEBPS/Date.xhtml",
		"/OEBPS/DateTime.xhtml",
		"/OEBPS/Dict.xhtml",
		"/OEBPS/Duration.xhtml",
		"/OEBPS/DynamicSupervisor.xhtml",
		"/OEBPS/Enum.EmptyError.xhtml",
		"/OEBPS/Enum.OutOfBoundsError.xhtml",
		"/OEBPS/Enum.xhtml",
		"/OEBPS/Enumerable.xhtml",
		"/OEBPS/ErlangError.xhtml",
		"/OEBPS/Exception.xhtml",
		"/OEBPS/File.CopyError.xhtml",
		"/OEBPS/File.Error.xhtml",
		"/OEBPS/File.LinkError.xhtml",
		"/OEBPS/File.RenameError.xhtml",
		"/OEBPS/File.Stat.xhtml",
		"/OEBPS/File.Stream.xhtml",
		"/OEBPS/File.xhtml",
		"/OEBPS/Float.xhtml",
		"/OEBPS/Function.xhtml",
		"/OEBPS/FunctionClauseError.xhtml",
		"/OEBPS/GenEvent.xhtml",
		"/OEBPS/GenServer.xhtml",
		"/OEBPS/HashDict.xhtml",
		"/OEBPS/HashSet.xhtml",
		"/OEBPS/IO.ANSI.xhtml",
		"/OEBPS/IO.Stream.xhtml",
		"/OEBPS/IO.StreamError.xhtml",
		"/OEBPS/IO.xhtml",
		"/OEBPS/Inspect.Algebra.xhtml",
		"/OEBPS/Inspect.Error.xhtml",
		"/OEBPS/Inspect.Opts.xhtml",
		"/OEBPS/Inspect.xhtml",
		"/OEBPS/Integer.xhtml",
		"/OEBPS/JSON.DecodeError.xhtml",
		"/OEBPS/JSON.Encoder.xhtml",
		"/OEBPS/JSON.xhtml",
		"/OEBPS/Kernel.ParallelCompiler.xhtml",
		"/OEBPS/Kernel.SpecialForms.xhtml",
		"/OEBPS/Kernel.TypespecError.xhtml",
		"/OEBPS/Kernel.xhtml",
		"/OEBPS/KeyError.xhtml",
		"/OEBPS/Keyword.xhtml",
		"/OEBPS/List.Chars.xhtml",
		"/OEBPS/List.xhtml",
		"/OEBPS/Macro.Env.xhtml",
		"/OEBPS/Macro.xhtml",
		"/OEBPS/Map.xhtml",
		"/OEBPS/MapSet.xhtml",
		"/OEBPS/MatchError.xhtml",
		"/OEBPS/MismatchedDelimiterError.xhtml",
		"/OEBPS/MissingApplicationsError.xhtml",
		"/OEBPS/Module.xhtml",
		"/OEBPS/NaiveDateTime.xhtml",
		"/OEBPS/Node.xhtml",
		"/OEBPS/OptionParser.ParseError.xhtml",
		"/OEBPS/OptionParser.xhtml",
		"/OEBPS/PartitionSupervisor.xhtml",
		"/OEBPS/Path.xhtml",
		"/OEBPS/Port.xhtml",
		"/OEBPS/Process.xhtml",
		"/OEBPS/Protocol.UndefinedError.xhtml",
		"/OEBPS/Protocol.xhtml",
		"/OEBPS/Range.xhtml",
		"/OEBPS/Record.xhtml",
		"/OEBPS/Regex.CompileError.xhtml",
		"/OEBPS/Regex.xhtml",
		"/OEBPS/Registry.xhtml",
		"/OEBPS/RuntimeError.xhtml",
		"/OEBPS/Set.xhtml",
		"/OEBPS/Stream.xhtml",
		"/OEBPS/String.Chars.xhtml",
		"/OEBPS/String.xhtml",
		"/OEBPS/StringIO.xhtml",
		"/OEBPS/Supervisor.Spec.xhtml",
		"/OEBPS/Supervisor.xhtml",
		"/OEBPS/SyntaxError.xhtml",
		"/OEBPS/System.EnvError.xhtml",
		"/OEBPS/System.xhtml",
		"/OEBPS/SystemLimitError.xhtml",
		"/OEBPS/Task.Supervisor.xhtml",
		"/OEBPS/Task.xhtml",
		"/OEBPS/Time.xhtml",
		"/OEBPS/TokenMissingError.xhtml",
		"/OEBPS/TryClauseError.xhtml",
		"/OEBPS/Tuple.xhtml",
		"/OEBPS/URI.Error.xhtml",
		"/OEBPS/URI.xhtml",
		"/OEBPS/UndefinedFunctionError.xhtml",
		"/OEBPS/UnicodeConversionError.xhtml",
		"/OEBPS/Version.InvalidRequirementError.xhtml",
		"/OEBPS/Version.InvalidVersionError.xhtml",
		"/OEBPS/Version.Requirement.xhtml",
		"/OEBPS/Version.xhtml",
		"/OEBPS/WithClauseError.xhtml",
		"/OEBPS/agents.xhtml",
		"/OEBPS/alias-require-and-import.xhtml",
		"/OEBPS/anonymous-functions.xhtml",
		"/OEBPS/assets/kv-observer.png",
		"/OEBPS/assets/logo.png",
		"/OEBPS/basic-types.xhtml",
		"/OEBPS/binaries-strings-and-charlists.xhtml",
		"/OEBPS/case-cond-and-if.xhtml",
		"/OEBPS/changelog.xhtml",
		"/OEBPS/code-anti-patterns.xhtml",
		"/OEBPS/compatibility-and-deprecations.xhtml",
		"/OEBPS/comprehensions.xhtml",
		"/OEBPS/config-and-releases.xhtml",
		"/OEBPS/content.opf",
		"/OEBPS/debugging.xhtml",
		"/OEBPS/dependencies-and-umbrella-projects.xhtml",
		"/OEBPS/design-anti-patterns.xhtml",
		"/OEBPS/dist/epub-LSJCIYTM.js",
		"/OEBPS/dist/epub-elixir-B7F5ZCEW.css",
		"/OEBPS/distributed-tasks.xhtml",
		"/OEBPS/docs-tests-and-with.xhtml",
		"/OEBPS/domain-specific-languages.xhtml",
		"/OEBPS/dynamic-supervisor.xhtml",
		"/OEBPS/enum-cheat.xhtml",
		"/OEBPS/enumerable-and-streams.xhtml",
		"/OEBPS/erlang-libraries.xhtml",
		"/OEBPS/erlang-term-storage.xhtml",
		"/OEBPS/genservers.xhtml",
		"/OEBPS/gradual-set-theoretic-types.xhtml",
		"/OEBPS/introduction-to-mix.xhtml",
		"/OEBPS/introduction.xhtml",
		"/OEBPS/io-and-the-file-system.xhtml",
		"/OEBPS/keywords-and-maps.xhtml",
		"/OEBPS/library-guidelines.xhtml",
		"/OEBPS/lists-and-tuples.xhtml",
		"/OEBPS/macro-anti-patterns.xhtml",
		"/OEBPS/macros.xhtml",
		"/OEBPS/module-attributes.xhtml",
		"/OEBPS/modules-and-functions.xhtml",
		"/OEBPS/naming-conventions.xhtml",
		"/OEBPS/nav.xhtml",
		"/OEBPS/operators.xhtml",
		"/OEBPS/optional-syntax.xhtml",
		"/OEBPS/pattern-matching.xhtml",
		"/OEBPS/patterns-and-guards.xhtml",
		"/OEBPS/process-anti-patterns.xhtml",
		"/OEBPS/processes.xhtml",
		"/OEBPS/protocols.xhtml",
		"/OEBPS/quote-and-unquote.xhtml",
		"/OEBPS/recursion.xhtml",
		"/OEBPS/sigils.xhtml",
		"/OEBPS/structs.xhtml",
		"/OEBPS/supervisor-and-application.xhtml",
		"/OEBPS/syntax-reference.xhtml",
		"/OEBPS/task-and-gen-tcp.xhtml",
		"/OEBPS/title.xhtml",
		"/OEBPS/try-catch-and-rescue.xhtml",
		"/OEBPS/typespecs.xhtml",
		"/OEBPS/unicode-syntax.xhtml",
		"/OEBPS/what-anti-patterns.xhtml",
		"/OEBPS/writing-documentation.xhtml",
		"/mimetype",
	}
	root := manifests.GenerateDirectoryTree(files)
	buf := make([]byte, 10000)
	out := bytes.NewBuffer(buf)
	writ := bufio.NewWriter(out)
	manifests.MarshalToc(root, writ)
	writ.Flush()
	println(string(out.Bytes()))
}
