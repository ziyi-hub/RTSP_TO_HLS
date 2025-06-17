const express = require('express')
const app = express()
const port = 3000

let cameraData = {
    "resultCode": 0,
    "cameraBriefInfos": {
    "total": 986,
        "indexRange": {
        "fromIndex": 1,
            "toIndex": 2000
    },
    "cameraBriefInfoList": {
        "cameraBriefInfo": [
            {
                "code": "42011201001310948425#17#4c6a6bade5c74d66b71971f6f2670e61",
                "name": "12-1#支洞-全景【球机】",
                "deviceGroupCode": "42011200002000202404#4c6a6bade5c74d66b71971f6f2670e61",
                "parentCode": "",
                "domainCode": "4c6a6bade5c74d66b71971f6f2670e61",
                "deviceModelType": "",
                "vendorType": "",
                "deviceFormType": 1,
                "type": 2,
                "cameraLocation": "42011201001310948425",
                "cameraStatus": 1,
                "status": 1,
                "netType": 0,
                "isSupportIntelligent": 0,
                "enableVoice": 1,
                "nvrCode": "",
                "deviceCreateTime": "",
                "isExDomain": 1,
                "deviceIP": "",
                "reserve": null
            },
            {
                "code": "42011201001310631072#17#4c6a6bade5c74d66b71971f6f2670e61",
                "name": "12#-支洞-全景【枪机】",
                "deviceGroupCode": "42011200002000202404#4c6a6bade5c74d66b71971f6f2670e61",
                "parentCode": "",
                "domainCode": "4c6a6bade5c74d66b71971f6f2670e61",
                "deviceModelType": "",
                "vendorType": "",
                "deviceFormType": 1,
                "type": 0,
                "cameraLocation": "42011201001310631072",
                "cameraStatus": 1,
                "status": 1,
                "netType": 0,
                "isSupportIntelligent": 0,
                "enableVoice": 1,
                "nvrCode": "",
                "deviceCreateTime": "",
                "isExDomain": 1,
                "deviceIP": "",
                "reserve": null
            },
        ]
    }
},
    "audioBriefInfos": null,
    "alarmBriefInfos": null,
    "cameraBriefInfosV2": null,
    "shadowCameraBriefInfos": null,
    "cameraBriefExInfos": null
}


let json =
    {
        "type": "collection", "count": 3, "commandes": [
            {
                "rtspURL": "rtsp://10.70.37.12:1166/42011201001310631072?Short=1&Token=f3T602CgJnCSkJqjZqnX2PfL6Dj2r7zsLbeo1zCgVUg=&DomainCode=4c6a6bade5c74d66b71971f6f2670e61&UserId=6&",
            },
            {
                "rtspURL": "rtsp://10.70.37.12:1166/41010400001310001290?Short=1&Token=u2TQrt0bu/BncUqDPE8bZNdSL7BDyaAl21Nq/y2SA4k=&DomainCode=4c6a6bade5c74d66b71971f6f2670e61&UserId=6&",
            },
            {
                "rtspURL": "rtsp://10.70.37.12:1166/11010000001320000050?Short=1&Token=XKegRFG839E4E0NR50fW49kPqj4idqUJT4opMc0tCVM=&DomainCode=4c6a6bade5c74d66b71971f6f2670e61&UserId=6&",
            }
        ],
    }

app.use(express.json()) // for parsing application/json
app.use(express.urlencoded({ extended: true })) // for parsing application/x-www-form-urlencoded

app.get("/cameras", async (req, res) => {
    res.json(cameraData.cameraBriefInfos.cameraBriefInfoList.cameraBriefInfo)
})

app.post("/video", async (req, res) => {
    const { cameraCode } = req.body;

    if (!cameraCode) {
        return res.status(400).json({ error: "Missing cameraCode" });
    }

    // 拆解 cameraCode
    const parts = cameraCode.split("#");
    if (parts.length !== 3) {
        return res.status(400).json({ error: "Invalid cameraCode format" });
    }

    const cameraId = parts[0];
    const domainCode = parts[2];

    const result = json.commandes.find(commande => {
        const url = commande.rtspURL;
        return (
            url.includes(cameraId) &&
            url.includes(`DomainCode=${domainCode}`)
        );
    });

    if (result) {
        res.json({ rtspURL: result.rtspURL, resultCode: 0 });
    } else {
        res.status(404).json({ error: "Video stream not found for provided cameraCode" });
    }
});

//
// app.post('/login', function (req, res){
//     const username = req.body.username
//     const user = {name: username}
//
//     // sign with RSA SHA256
//     const accessToken = jwt.sign({user}, 'my_secret_key', {algorithm: 'HS256', expiresIn: '600s'})
//     return res.status(201).json({
//         user: req.headers.nom,
//         token: accessToken
//     });
// })

// function auth(req, res, next) {
//     if(req.headers['authorization']){
//         const authHeader = req.headers['authorization']
//         const token = authHeader && authHeader.split(' ')[1]
//         if(token !== null){
//             jwt.verify(token, 'my_secret_key', (err, user) => {
//                 if(err){
//                     res.status(500).json({
//                         error: "Not Authorized"
//                     });
//                     throw new Error("Not Authorized");
//                 }
//                 req.user = user
//                 next();
//             })
//         }
//     }else{
//         return res.status(401).json({
//             error: "Missing credential"
//         })
//     }
// }

app.use((function (req, res) {
    res.sendStatus(404)
}))

app.listen(port, () => {
    console.log(`Example app listening at http://localhost:${port}`)
})




