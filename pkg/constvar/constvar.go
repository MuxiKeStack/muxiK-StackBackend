package constvar

const (
	DefaultLimit          = 50
	AllowRemindAdminLimit = 5
	DefaultAdminEmailAddr = "1027319981@qq.com"
	EmailTemp             = `<table border="0" cellspacing="0" cellpadding="0" style="max-width: 600px;">
	<tbody>
		<tr height="16"></tr>
		<tr>
			<td>
				<table bgcolor="#4184F3" width="100%" border="0" cellspacing="0" cellpadding="0" style="min-width: 332px; max-width: 600px; border: 1px solid #E0E0E0; border-bottom: 0; border-top-left-radius: 3px; border-top-right-radius: 3px;">
					<tbody>
						<tr>
							<td height="48px" colspan="3"></td>
						</tr>
						<tr>
							<td width="32px"></td>
							<td style="font-family: Roboto-Regular,Helvetica,Arial,sans-serif; font-size: 24px; color: #FFFFFF; line-height: 1.25;">木犀课栈: 评课 举报通知</td>
							<td width="32px"></td>
						</tr>
						<tr>
							<td height="18px" colspan="3"></td>
						</tr>
					</tbody>
				</table>
			</td>
		</tr>
		<tr>
			<td>
				<table bgcolor="#FAFAFA" width="100%" border="0" cellspacing="0" cellpadding="0" style="min-width: 332px; max-width: 600px; border: 1px solid #F0F0F0; border-bottom: 1px solid #C0C0C0; border-top: 0; border-bottom-left-radius: 3px; border-bottom-right-radius: 3px;">
					<tbody>
						<tr height="16px">
							<td width="32px" rowspan="3"></td>
							<td></td>
							<td width="32px" rowspan="3"></td>
						</tr>
						<tr>
							<td>
								<p>尊敬的 木犀课栈 掌柜的：</p>
								<p>驿站检测到了多封对于 评课（ID：EVALUATION_ID）的举报信，该评课被举报次数为：</p>
									<div style="text-align: center;">
										<p dir="ltr">
											<strong style="text-align: kecenter; font-size: 24px; font-weight: bold;">REPORT_TOT</strong>
										</p>
									</div>
									<p>以下是所有的举报该 评课 的来自五湖四海的内容：</p>
									ALL_REASON
									<p>此致</p>
									<p>MuxiStudio 木犀课栈敬上</p>
								</td>
							</tr>
							<tr height="32px"></tr>
						</tbody>
					</table>
				</td>
			</tr>
			<tr height="16"></tr>
			<tr>
				<td style="max-width: 600px; font-family: Roboto-Regular,Helvetica,Arial,sans-serif; font-size: 10px; color: #BCBCBC; line-height: 1.5;">
				</td>
			</tr>
			<tr>
				<td>此电子邮件地址无法接收回复。如需更多信息，请访问 
					<a href="https://www.muxixyz.com" style="text-decoration: none; color: #4d90fe;" rel="noopener" target="_blank">木犀团队官网
					</a>。
					<br>Muxi Studio, CCNU, Wuhan HuBei, China 
					<table style="font-family: Roboto-Regular,Helvetica,Arial,sans-serif; font-size: 10px; color: #666666; line-height: 18px; padding-bottom: 10px">						
					</table>
				</td>
			</tr>
		</tbody>
	</table>`
)

var (
	Attendance map[uint8]string = map[uint8]string{
		1: "经常点名",
		2: "偶尔点名",
		3: "签到点名",
	}
	Exam map[uint8]string = map[uint8]string{
		1: "无考核",
		2: "闭卷考试",
		3: "开卷考试",
		4: "论文考核",
	}
	AttendanceEnglish map[uint8]string = map[uint8]string{
		1: "Often",
		2: "Occasionally",
		3: "SignIn",
	}
	ExamEnglish map[uint8]string = map[uint8]string{
		1: "None",
		2: "Close",
		3: "Open",
		4: "Eassay",
	}
)
