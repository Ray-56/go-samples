INSERT INTO postbrother_4 (callee, bizlabel, starttime, endtime) values ('-9501801111222','crd1-aaabbb', '2022-04-04T09:42:42+08:00', '2022-04-04T09:42:44+08:00');
INSERT INTO postbrother_5 (callee, bizlabel, starttime, endtime) values ('+9501310000111','crd2-biz2', '2022-05-04T09:42:42+08:00', '2022-05-04T09:42:44+08:00');
INSERT INTO postbrother_6 (callee, bizlabel, starttime, endtime) values ('+9501310000111','这里是第2张表的公司名222222', '2022-06-06T09:42:42+08:00', '2022-06-6T09:42:44+08:00');



SELECT * FROM postbrother_4
WHERE callee LIKE '%+95%';
