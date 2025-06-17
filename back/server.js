const express = require('express')
const app = express()
const port = 3000

let cameras = [
    "41010400001310001212#16#4c6a6bade5c74d66b71971f6f2670e61",
    "41010400001310001290#16#4c6a6bade5c74d66b71971f6f2670e61",
    "11010000001320000050#16#4c6a6bade5c74d66b71971f6f2670e61"
]

let json =
    {
        "type": "collection", "count": 3, "commandes": [
            {
                "rtspURL": "rtsp://10.70.37.12:1166/41010400001310001212?Short=1&Token=f3T602CgJnCSkJqjZqnX2PfL6Dj2r7zsLbeo1zCgVUg=&DomainCode=4c6a6bade5c74d66b71971f6f2670e61&UserId=6&",
            },
            {
                "rtspURL": "rtsp://10.70.37.12:1166/41010400001310001290?Short=1&Token=u2TQrt0bu/BncUqDPE8bZNdSL7BDyaAl21Nq/y2SA4k=&DomainCode=4c6a6bade5c74d66b71971f6f2670e61&UserId=6&",
                // "rtspURL2": "rtsp://stream.strba.sk:1935/strba/VYHLAD_JAZERO.stream"
            },
            {
                "rtspURL": "rtsp://10.70.37.12:1166/11010000001320000050?Short=1&Token=XKegRFG839E4E0NR50fW49kPqj4idqUJT4opMc0tCVM=&DomainCode=4c6a6bade5c74d66b71971f6f2670e61&UserId=6&",
            }
        ],
    }

app.use(express.json()) // for parsing application/json
app.use(express.urlencoded({ extended: true })) // for parsing application/x-www-form-urlencoded

app.get("/cameras", async(req, res, next) => {
    res.json(cameras)
})

app.post("/video", async (req, res) => {
    const { cameraCode, mediaURLParam } = req.body;

    if (!cameraCode || !mediaURLParam) {
        return res.status(400).json({ error: "Missing cameraCode or mediaURLParam" });
    }

    // 拆解 cameraCode
    const parts = cameraCode.split("#");
    if (parts.length !== 3) {
        return res.status(400).json({ error: "Invalid cameraCode format" });
    }

    const cameraId = parts[0]; // e.g., "41010400001310001290"
    const domainCode = parts[2]; // e.g., "4c6a6bade5c74d66b71971f6f2670e61"

    // 在 json.commandes 中查找匹配的 URL
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




