package vars

const (
	FileHelpEn = "Specifies the input file path."
	FileHelpZh = "指定输入文件路径。"

	OutputHelpEn = "Specifies the output file path."
	OutputHelpZh = "指定输出文件路径。"

	IPHelpEn = "Matches IPs from the input pipe or file."
	IPHelpZh = "从输入管道或文件中匹配 IP。"

	ExcludeHelpEn = "Excludes internal/private IP segments when using -i/--ip."
	ExcludeHelpZh = "在使用 -i/--ip 时排除内部/私有 IP 段。"

	DomainHelpEn = "Matches domains from the input pipe or file."
	DomainHelpZh = "从输入管道或文件中匹配域名。"

	RootDomainHelpEn = "Outputs only the primary domain when using -d/--domain."
	RootDomainHelpZh = "在使用 -d/--domain 时仅输出主要域名。"

	WithPortHelpEn = "Filters only domain & IP:port combinations."
	WithPortHelpZh = "仅筛选域名和 IP:端口 组合。"

	RuleHelpEn = "Utilizes a custom replacement rule (custom output replacement rule: https://{}/)."
	RuleHelpZh = "使用自定义输出替换规则（自定义输出替换规则：https://{}/）。"

	FlagHelpEn = "Specifies the replacement identification."
	FlagHelpZh = "指定替换标识。"

	URLHelpEn = "Matches URLs from the input pipe or file."
	URLHelpZh = "从输入管道或文件中匹配 URL。"

	URLFilterHelpEn = "Filters URLs with specific extensions."
	URLFilterHelpZh = "使用特定扩展名过滤 URL。"

	CidrHelpEn = "Outputs the specified CIDR IP list."
	CidrHelpZh = "输出指定 CIDR 范围内的所有 IP。"

	LimitLenHelpEn = "Matches input specified length string, e.g., \"-l 35\" == \"-l 0-35\"."
	LimitLenHelpZh = "匹配每行指定长度的字符串，例如，\"-l 35\" == \"-l 0-35\"。"

	ShowHelpEn = "Displays the length of each line and provides summaries."
	ShowHelpZh = "显示每行的长度并提供摘要。"

	ProgressHelpEn = "Outputs execution progress metrics."
	ProgressHelpZh = "读取大量行时输出执行进度指标状态。"

	UpdateHelpEn = "Updates the tool engine to the latest released version."
	UpdateHelpZh = "将工具引擎更新到最新版本。"

	GrepPatternHelpEn = "Pattern for regex."
	GrepPatternHelpZh = "正则表达式模式。"

	InverseMatchHelpEn = "Invert the match pattern."
	InverseMatchHelpZh = "反转匹配模式。"

	DiffCmdHelpEn = "Compares files using different modes:\n1: A-B\n2: B-A\n3: A&B"
	DiffCmdHelpZh = "使用不同模式比较文件：\n1：A-B\n2：B-A\n3：A&B"

	StrictModeHelpEn = "Match lines strictly one by one (non-default)."
	StrictModeHelpZh = "严格逐行匹配（非默认）。"

	SmartHelpEn = "Use heuristic technique to remove duplicated lines."
	SmartHelpZh = "使用启发式技术去除重复行。"

	ThresholdHelpEn = "Set threshold for smart strategy."
	ThresholdHelpZh = "设置智能策略的阈值。"

	AlterHelpEn = "IP Alters (0,1,2,3,4,5,6,7,8)"
	AlterHelpZh = "IP 变换 (0,1,2,3,4,5,6,7,8)"
)
