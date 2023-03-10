<#escape x as x?xml>
<w:document xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006" xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml" xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml" xmlns:w16="http://schemas.microsoft.com/office/word/2018/wordml" xmlns:w16cex="http://schemas.microsoft.com/office/word/2018/wordml/cex" xmlns:w16cid="http://schemas.microsoft.com/office/word/2016/wordml/cid" xmlns:w16se="http://schemas.microsoft.com/office/word/2015/wordml/symex" xmlns:wp14="http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing" mc:Ignorable="w14 w15 w16se w16cid w16 w16cex wp14">
    <w:body>
        <w:tbl>
            <w:tblPr>
                <w:tblStyle w:val="ab"/>
                <w:tblW w:type="auto" w:w="0"/>
                <w:tblInd w:type="dxa" w:w="137"/>
                <w:tblLook w:firstColumn="1" w:firstRow="1" w:lastColumn="0" w:lastRow="0" w:noHBand="0" w:noVBand="1" w:val="04A0"/>
            </w:tblPr>
            <w:tblGrid>
                <w:gridCol w:w="8493"/>
            </w:tblGrid>
            <#list table as row>
            <w:tr w:rsidR="00DF50D4" w:rsidTr="006F0787">
                <w:trPr>
                    <w:trHeight w:val="462"/>
                </w:trPr>
                <w:tc>
                    <w:tcPr>
                        <w:tcW w:type="dxa" w:w="9043"/>
                    </w:tcPr>
                    <w:p w:rsidP="006F0787" w:rsidR="00DF50D4" w:rsidRDefault="00DF50D4">
                        <w:pPr>
                            <w:widowControl/>
                            <w:spacing w:after="0" w:line="240" w:lineRule="auto"/>
                            <w:jc w:val="center"/>
                        </w:pPr>
                        <w:r w:rsidRPr="008A79E2">
                            <w:rPr>
                                <w:rFonts w:ascii="微软雅黑" w:eastAsia="微软雅黑" w:hAnsi="微软雅黑"/>
                                <w:sz w:val="21"/>
                                <w:szCs w:val="21"/>
                                <w:lang w:eastAsia="zh-CN"/>
                            </w:rPr>
                            <w:t>${row}</w:t>
                        </w:r>
                    </w:p>
                </w:tc>
            </w:tr>
            </#list>
        </w:tbl>
        <w:p w:rsidR="00000000" w:rsidRDefault="00DF50D4"/>
    </w:body>
</w:document>
</#escape>