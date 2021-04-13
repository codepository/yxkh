# 同步至fznews_flow_process表

0. 将hrm数据库中: process, info_user,res_evaluation,res_responsibility,res_project,res_mark 表 传输至 fzrbwx
1. fzrbwx数据库: 清空 fznews_flow_process 表 

# 同步process表
0. 将hrm数据库中: process, info_user 表 传输至 fzrbwx
1. 然后按顺序执行下列命令
INSERT INTO fznews_flow_process(step,deptName,completed,businessType,processInstanceId,uid,requestedDate,title)SELECT 3,deptName,1,businessType,processInstanceId,userId,requestedDate,title from process;
update fznews_flow_process set businessType=SUBSTR(businessType,6) where businessType like '一线干部-%';
update fznews_flow_process set step=3;
UPDATE fznews_flow_process f,fznews_userinfo_mapping m set f.username=m.name,f.userId=m.userid,f.uid=m.id where f.uid=m.uid;
update fznews_flow_process set businessType=SUBSTR(businessType,5) where businessType like '一线干部%';

2. 删除fzrbwx数据的process表

# 同步res_evaluation表
将表 res_evaluation 复制到fzrbwx数据库，删除外键

将 res_responsibility 复制到fzrbwx数据库

update res_responsibility set sparation=replace(sparation,"一线干部","-");
INSERT INTO res_evaluation(startDate,endDate,processInstanceId,selfEvaluation,shortComesAndPlan,sparation,createTime) SELECT startDate,endDate,processInstanceId,currentJob,currentJobDescription,sparation,createTime from res_responsibility;
update res_evaluation r,fznews_flow_process p set r.uid=p.uid,r.username=p.username,r.department=p.deptName where r.processInstanceId=p.processInstanceId;
update res_evaluation r,weixin_leave_userinfo u set r.position=u.position where r.uid=u.id;

最后将表剪切到fznews_yxkh数据库

# 同步res_project

将 res_project 复制到fzrbwx数据库

update res_project r,fznews_userinfo_mapping m set r.userId=m.id where r.userId=m.uid;



# 同步res_mark
将 res_mark 复制到fzrbwx数据库

update res_mark r,fznews_userinfo_mapping m set r.userId=m.id,r.username=m.name where r.userId=m.uid;



最后将表: res_project,res_mark,res_evaluation,剪切到fznews_yxkh数据库
