<w:hdr xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas"
       xmlns:cx="http://schemas.microsoft.com/office/drawing/2014/chartex"
       xmlns:cx1="http://schemas.microsoft.com/office/drawing/2015/9/8/chartex"
       xmlns:cx2="http://schemas.microsoft.com/office/drawing/2015/10/21/chartex"
       xmlns:cx3="http://schemas.microsoft.com/office/drawing/2016/5/9/chartex"
       xmlns:cx4="http://schemas.microsoft.com/office/drawing/2016/5/10/chartex"
       xmlns:cx5="http://schemas.microsoft.com/office/drawing/2016/5/11/chartex"
       xmlns:cx6="http://schemas.microsoft.com/office/drawing/2016/5/12/chartex"
       xmlns:cx7="http://schemas.microsoft.com/office/drawing/2016/5/13/chartex"
       xmlns:cx8="http://schemas.microsoft.com/office/drawing/2016/5/14/chartex"
       xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"
       xmlns:aink="http://schemas.microsoft.com/office/drawing/2016/ink"
       xmlns:am3d="http://schemas.microsoft.com/office/drawing/2017/model3d"
       xmlns:o="urn:schemas-microsoft-com:office:office"
       xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
       xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math"
       xmlns:v="urn:schemas-microsoft-com:vml"
       xmlns:wp14="http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing"
       xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"
       xmlns:w10="urn:schemas-microsoft-com:office:word"
       xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
       xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml"
       xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml"
       xmlns:w16cex="http://schemas.microsoft.com/office/word/2018/wordml/cex"
       xmlns:w16cid="http://schemas.microsoft.com/office/word/2016/wordml/cid"
       xmlns:w16="http://schemas.microsoft.com/office/word/2018/wordml"
       xmlns:w16se="http://schemas.microsoft.com/office/word/2015/wordml/symex"
       xmlns:wpg="http://schemas.microsoft.com/office/word/2010/wordprocessingGroup"
       xmlns:wpi="http://schemas.microsoft.com/office/word/2010/wordprocessingInk"
       xmlns:wne="http://schemas.microsoft.com/office/word/2006/wordml"
       xmlns:wps="http://schemas.microsoft.com/office/word/2010/wordprocessingShape"
       mc:Ignorable="w14 w15 w16se w16cid w16 w16cex wp14">
    <w:p w14:paraId="576D8AF3" w14:textId="302AE919" w:rsidR="00E749B5" w:rsidRDefault="00094624"
         w:rsidP="00283429">
        <w:pPr>
            <w:pBdr>
                <w:bottom w:val="single" w:sz="2" w:space="1" w:color="auto"/>
            </w:pBdr>
            <w:ind w:right="180"/>
            <w:jc w:val="right"/>
            <w:rPr>
                <w:rFonts w:ascii="微软雅黑" w:eastAsia="微软雅黑" w:hAnsi="微软雅黑"/>
                <w:sz w:val="18"/>
                <w:szCs w:val="18"/>
                <w:lang w:eastAsia="zh-CN"/>
            </w:rPr>
        </w:pPr>
        <#if watermarkRelationshipId??>
            <w:r>
                <w:rPr>
                    <w:noProof/>
                </w:rPr>
                <w:pict w14:anchorId="793D3C32">
                    <v:shapetype id="_x0000_t75" coordsize="21600,21600" o:spt="75" o:preferrelative="t"
                                 path="m@4@5l@4@11@9@11@9@5xe" filled="f" stroked="f">
                        <v:stroke joinstyle="miter"/>
                        <v:formulas>
                            <v:f eqn="if lineDrawn pixelLineWidth 0"/>
                            <v:f eqn="sum @0 1 0"/>
                            <v:f eqn="sum 0 0 @1"/>
                            <v:f eqn="prod @2 1 2"/>
                            <v:f eqn="prod @3 21600 pixelWidth"/>
                            <v:f eqn="prod @3 21600 pixelHeight"/>
                            <v:f eqn="sum @0 0 1"/>
                            <v:f eqn="prod @6 1 2"/>
                            <v:f eqn="prod @7 21600 pixelWidth"/>
                            <v:f eqn="sum @8 21600 0"/>
                            <v:f eqn="prod @7 21600 pixelHeight"/>
                            <v:f eqn="sum @10 21600 0"/>
                        </v:formulas>
                        <v:path o:extrusionok="f" gradientshapeok="t" o:connecttype="rect"/>
                        <o:lock v:ext="edit" aspectratio="t"/>
                    </v:shapetype>
                    <v:shape id="WordPictureWatermark42581196" o:spid="_x0000_s2049" type="#_x0000_t75" alt=""
                             style="position:absolute;margin-left:0;margin-top:0;width:461.3pt;height:193.55pt;z-index:-251656192;mso-wrap-edited:f;mso-width-percent:0;mso-height-percent:0;mso-position-horizontal:center;mso-position-horizontal-relative:margin;mso-position-vertical:center;mso-position-vertical-relative:margin;mso-width-percent:0;mso-height-percent:0"
                             o:allowincell="f">
                        <v:imagedata r:id="${watermarkRelationshipId?xml}" o:title="logo" gain="19661f" blacklevel="22938f"/>
                        <w10:wrap anchorx="margin" anchory="margin"/>
                    </v:shape>
                </w:pict>
            </w:r>
        </#if>
        <w:r w:rsidR="00E749B5" w:rsidRPr="00E95C2C">
            <w:rPr>
                <w:rFonts w:ascii="微软雅黑" w:eastAsia="微软雅黑" w:hAnsi="微软雅黑"/>
                <w:sz w:val="18"/>
                <w:szCs w:val="18"/>
                <w:lang w:eastAsia="zh-CN"/>
            </w:rPr>
            <w:t>${header_text}</w:t>
        </w:r>
    </w:p>
</w:hdr>
