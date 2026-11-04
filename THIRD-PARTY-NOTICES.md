# Third-party notices

My Voice source is MIT. Included support components retain their own licences.
No model weights or personal reference recordings are bundled in the executable.

| Included component | Version | Licence / notice |
|---|---|---|
| Go runtime and server modules | Go1.27.3; exact modules in docs/go-dependencies.json and go.sum | BSD/MIT/Apache notices in licenses/go and licenses/Go-BSD.txt |
| One Workspace apphost | v0.3.0-preview.110 | MIT; oneworkspacestudio/one-workspace |
| React / React DOM / scheduler | 19.3.1 / 19.3.1 / 0.28.0 | MIT; licenses/web |
| Lucide | 0.577.0 | ISC; licenses/web |
| Codedia AUI compiler | 0.2.1, build tool | MIT; licenses/web |
| Pocket native Rust CPU runtime | official v4.0.1 archive, CLI3.3.0; baseline x86-64 build of the same pinned source | MIT; licenses/Pocket-MIT.txt |
| MOSS CPU runtime source | 8b7bcc9341b3b4ef3a3a58ba1338a7d85ff133eb | Apache-2.0; licenses/MOSS-Apache-2.0.txt |
| CPython embedded x64 | 3.12.10 | PSF and incorporated notices; licenses/Python-3.12.txt |
| NumPy | 2.4.4 | BSD and incorporated numerical-library notices; licenses/python |
| ONNX Runtime | 1.24.4, CPU only | MIT and incorporated notices; licenses/python/onnxruntime |
| SentencePiece | 0.2.1 | Apache-2.0 |
| Flatbuffers / packaging / protobuf / SymPy / mpmath | exact pinned wheels in runtime.lock.json | Apache-2.0 / BSD / BSD / BSD / BSD, upstream wheel notices |
| FFmpeg audio runtime | 9.0.2, commit1e65c6cd5d | LGPL-2.1-or-later; licenses/FFmpeg-LGPL-2.1.txt |
| LAME | 3.100, statically included in FFmpeg | LGPL-2.0-or-later; licenses/LAME-LGPL-2.0.txt |

The small FFmpeg executable is built natively with Zig0.15.2, without GPL or
nonfree configuration, network protocols or autodetected external codecs.
Its complete FFmpeg/LAME corresponding sources and exact build recipe are offered
as `audio-runtime-corresponding-source.zip` beside every Windows release. The
recipe is also scripts/build-audio-runtime.sh; source/archive digests are locked
in runtime.lock.json. Users may modify/rebuild this LGPL program and replace the
`ffmpeg.exe` member in the source build's payload; prepare-runtime.py and the Go
source permit relinking/rebuilding the launcher. No restriction on reverse
engineering these libraries for debugging modifications is imposed.

Python support uses prebuilt pinned wheels without PyTorch. The pinned wheel
archives remain available from the URLs in runtime.lock.json. Each wheel's own
licences and bundled dependency notices apply. Models and the CC0 Open Voice
reference download separately on explicit acceptance: see MODEL-LICENSES.md.

No ElevenLabs code, branding, credentials or proprietary models are included.


Pocket's runtime incorporates Rust dependencies, including Symphonia audio
components under MPL-2.0. `licenses/Pocket-THIRD-PARTY-NOTICES.txt` retains the
upstream copyright/licence texts; `docs/pocket-dependencies.json` records the
complete pinned upstream Cargo.lock inventory (including optional/platform
crates). Not every locked crate is compiled into the Windows executable.
`pocket-runtime-corresponding-source.zip` supplies the runtime source and locked
dependency source beside this release. Test audio and model files are excluded;
the original upstream URLs and archive digests are included. This is in addition
to the separate LGPL FFmpeg/LAME source archive. The MPL-covered files remain
available under their original licence; the enclosing My Voice code stays MIT.

The additional portable Pocket executable is built on native Windows/MSVC from
the same source commit and Cargo.lock. Its flags are `-C target-cpu=x86-64
-C target-feature=+crt-static`, overriding upstream's x86-64-v3 default.
`scripts/build-pocket-runtime.py` is the exact offline build recipe; its source
archive contains that recipe, the original upstream source and all locked crate
archives. The recipe only normalizes the Git dependency's workspace manifest
inheritance for Cargo's vendored source format; library code remains unchanged.


## Qwen native CPU runtime (preview8)

llama.cpp b12084, source f9448cddc48c129a4abc7ff9c2a84c2b4c101c61.
Official Windows CPU binaries and app-local Microsoft CRT are pinned in runtime.lock.json.
Only release CRT files are redistributed unmodified; no debug runtime.
Qwen weights remain separate opt-in downloads; see MODEL-LICENSES.md.

### llama.cpp MIT license

MIT License

Copyright (c) 2023-2026 The ggml authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.


### LLVM OpenMP license

==============================================================================
The LLVM Project is under the Apache License v2.0 with LLVM Exceptions:
==============================================================================

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

    TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

    1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

    2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

    3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

    4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

    5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

    6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

    7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

    8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

    9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

    END OF TERMS AND CONDITIONS

    APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

    Copyright [yyyy] [name of copyright owner]

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.


---- LLVM Exceptions to the Apache 2.0 License ----

As an exception, if, as a result of your compiling your source code, portions
of this Software are embedded into an Object form of such source code, you
may redistribute such embedded portions in such Object form without complying
with the conditions of Sections 4(a), 4(b) and 4(d) of the License.

In addition, if you combine or link compiled forms of this Software with
software that is licensed under the GPLv2 ("Combined Software") and if a
court of competent jurisdiction determines that the patent provision (Section
3), the indemnity provision (Section 9) or other Section of the License
conflicts with the conditions of the GPLv2, you may retroactively and
prospectively choose to deem waived or otherwise exclude such Section(s) of
the License, but only in their entirety and only with respect to the Combined
Software.

==============================================================================
Software from third parties included in the LLVM Project:
==============================================================================
The LLVM Project contains third party software which is under different license
terms. All such code will be identified clearly using at least one of two
mechanisms:
1) It will be in a separate directory tree with its own `LICENSE.txt` or
   `LICENSE` file at the top containing the specific license and restrictions
   which apply to that software, or
2) It will contain specific license and restriction terms at the top of every
   file.

==============================================================================
Legacy LLVM License (https://llvm.org/docs/DeveloperPolicy.html#legacy):
==============================================================================

The software contained in this directory tree is dual licensed under both the
University of Illinois "BSD-Like" license and the MIT license.  As a user of
this code you may choose to use it under either license.  As a contributor,
you agree to allow your code to be used under both.  The full text of the
relevant licenses is included below.

In addition, a license agreement from the copyright/patent holders of the
software contained in this directory tree is included below.

==============================================================================

University of Illinois/NCSA
Open Source License

Copyright (c) 1997-2019 Intel Corporation

All rights reserved.

Developed by:
    OpenMP Runtime Team
    Intel Corporation
    http://www.openmprtl.org

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal with
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

    * Redistributions of source code must retain the above copyright notice,
      this list of conditions and the following disclaimers.

    * Redistributions in binary form must reproduce the above copyright notice,
      this list of conditions and the following disclaimers in the
      documentation and/or other materials provided with the distribution.

    * Neither the names of Intel Corporation OpenMP Runtime Team nor the
      names of its contributors may be used to endorse or promote products
      derived from this Software without specific prior written permission.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
CONTRIBUTORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS WITH THE
SOFTWARE.

==============================================================================

Copyright (c) 1997-2019 Intel Corporation

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

==============================================================================

Intel Corporation

Software Grant License Agreement ("Agreement")

Except for the license granted herein to you, Intel Corporation ("Intel") reserves
all right, title, and interest in and to the Software (defined below).

Definition

"Software" means the code and documentation as well as any original work of
authorship, including any modifications or additions to an existing work, that
is intentionally submitted by Intel to llvm.org (http://llvm.org) ("LLVM") for
inclusion in, or documentation of, any of the products owned or managed by LLVM
(the "Work"). For the purposes of this definition, "submitted" means any form of
electronic, verbal, or written communication sent to LLVM or its
representatives, including but not limited to communication on electronic
mailing lists, source code control systems, and issue tracking systems that are
managed by, or on behalf of, LLVM for the purpose of discussing and improving
the Work, but excluding communication that is conspicuously marked otherwise.

1. Grant of Copyright License. Subject to the terms and conditions of this
   Agreement, Intel hereby grants to you and to recipients of the Software
   distributed by LLVM a perpetual, worldwide, non-exclusive, no-charge,
   royalty-free, irrevocable copyright license to reproduce, prepare derivative
   works of, publicly display, publicly perform, sublicense, and distribute the
   Software and such derivative works.

2. Grant of Patent License. Subject to the terms and conditions of this
   Agreement, Intel hereby grants you and to recipients of the Software
   distributed by LLVM a perpetual, worldwide, non-exclusive, no-charge,
   royalty-free, irrevocable (except as stated in this section) patent license
   to make, have made, use, offer to sell, sell, import, and otherwise transfer
   the Work, where such license applies only to those patent claims licensable
   by Intel that are necessarily infringed by Intel's Software alone or by
   combination of the Software with the Work to which such Software was
   submitted. If any entity institutes patent litigation against Intel or any
   other entity (including a cross-claim or counterclaim in a lawsuit) alleging
   that Intel's Software, or the Work to which Intel has contributed constitutes
   direct or contributory patent infringement, then any patent licenses granted
   to that entity under this Agreement for the Software or Work shall terminate
   as of the date such litigation is filed.

Unless required by applicable law or agreed to in writing, the software is
provided on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied, including, without limitation, any warranties or
conditions of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
PARTICULAR PURPOSE.

==============================================================================


### Microsoft Visual C++ V14 Redistributable and Runtime

Source: https://visualstudio.microsoft.com/wp-content/uploads/2025/10/Visual-C-V14-License-Redistributable_and_Runtime_ENU.docx

MICROSOFT SOFTWARE LICENSE TERMS
MICROSOFT VISUAL C++ V14 REDISTRIBUTABLE and RUNTIME 
Last Updated: October 1, 2025
These license terms are an agreement between Microsoft Corporation (or based on where you live, one of its affiliates) and you. They apply to the software named above. The terms also apply to any Microsoft services or updates for the software, except to the extent those have different terms.
Unless otherwise stated, these license terms remain in effect for this version of the software and any updates or revisions released. Microsoft may update or modify these license terms at any time. It is your responsibility to regularly review the license terms and any related announcements. The most current version of the license terms will always be available at https://aka.ms/VisualStudioLicenses. By continuing to use the software after changes are posted, you agree to the updated terms.
BY USING THE SOFTWARE, YOU ACCEPT THESE TERMS. IF YOU DO NOT ACCEPT THEM, DO NOT USE THE SOFTWARE.
IF YOU COMPLY WITH THESE LICENSE TERMS, YOU HAVE THE RIGHTS BELOW.
INSTALLATION AND USE RIGHTS. 
You may install and use any number of copies of the software.
FEEDBACK. If you give feedback, suggestions, or recommendations (collectively, “Feedback”) about the software to Microsoft, you give to Microsoft, without charge, the right to use, share, and commercialize your Feedback in any way and for any purpose. You will not give Feedback that is a third party’s confidential information or subject to a license that requires Microsoft to license its software or documentation to third parties because we include your Feedback in them. These rights survive this agreement.
TERMS FOR SPECIFIC COMPONENTS.
Third Party or Open Source Components. The software may include third party or open source components with separate legal notices or governed by other agreements, as may be described in the notices file(s) accompanying the software. 
SCOPE OF LICENSE. The software is licensed, not sold. These license terms only give you some rights to use the software. Microsoft reserves all other rights. Unless applicable law gives you more rights despite this limitation, you may use the software only as expressly permitted in these license terms. In doing so, you must comply with any technical limitations in the software that only allow you to use it in certain ways. For example, if Microsoft technically limits or disables extensibility for the software, you may not extend the software by, among other things, loading or injecting into the software any non-Microsoft add-ins, macros, or packages; modifying the software registry settings; or adding features or functionality equivalent to that found in Microsoft products and services. You may not
work around any technical limitations in the software;
reverse engineer, decompile or disassemble the software, or otherwise attempt to derive the source code for the software, except and only to the extent required by third party licensing terms governing use of certain open source components that may be included in the software;
remove, minimize, block or modify any notices of Microsoft or its suppliers in the software; 
use the software in any way that is against the law; 
share, publish, rent, lease, or otherwise distribute the software or any of its code; or 
provide the software as a stand-alone offering or combined with any of your applications for others to use, or transfer the software or these license terms to any third party.
COMPLIANCE WITH TRADE LAWS. Microsoft products, software, technology, and services (“Items”) may be subject to U.S. and other countries’ export jurisdictions. Each party will comply with all laws and regulations applicable to the import or export of the Items, including, without limitation, trade laws of the U.S., EU, and UK, such as the U.S. Export Administration Regulations, sanctions regulations administered by the U.S. Office of Foreign Assets Control, the EU Dual Use Regulation 2021/821, and/or other end-user, end use, and destination restrictions (“Trade Laws”) as well as the global legal compliance standards detailed in the Microsoft Standards of Business Conduct. You will not, and will ensure your Affiliates will not, take any action that causes Microsoft to violate applicable Trade Laws. Microsoft may suspend or terminate this agreement immediately without notice to the extent that Microsoft reasonably believes that performance would cause it to violate Trade Laws or put it at risk of becoming subject to sanctions and penalties under such laws. You remain responsible for its and for your Affiliates’ compliance with this section. For additional information, see www.microsoft.com/exporting.
You must ensure that you, any of your directors or officers, or to the best of your knowledge, your employees or agents, Affiliates, or Affiliates’ directors, officers, employees, or agents are not sanctioned persons designated under applicable Trade Laws.
You shall provide Microsoft with prompt written notice if you or any of your Affiliates receive notice from any government agency of any investigation, inquiry, or proceeding for any actual or alleged noncompliance of or with any Trade Laws related to activities subject to this Agreement.
You shall comply with global use restrictions in Microsoft’s product or service terms relevant to this agreement, as updated from time to time.
You shall not make available Microsoft Items to any entity or individual designated on any export or sanctions restriction lists under Trade Laws, in each case, without specific government authorization required by Trade Laws.
For the purposes of this Section 8, “Affiliate” means any legal entity that directly or indirectly owns, is owned by, or is commonly owned with the applicable party. “Own” means having more than 50% ownership or the right to direct the management of the entity.
SUPPORT. Because the software is “as is,” we may not provide support services for it.
ENTIRE AGREEMENT. These license terms, and the terms for supplements, updates, Internet-based services and support services, are the entire agreement for the software and support services.
APPLICABLE LAW. If you acquired the software in the United States, Washington state law applies to interpretation of and claims for breach of these license terms, and the laws of the state where you live apply to all other claims. If you acquired the software in any other country, its laws apply.
CONSUMER RIGHTS; REGIONAL VARIATIONS. These license terms describe certain legal rights. You may have other rights, including consumer rights, under the laws of your state or country. Separate and apart from your relationship with Microsoft, you may also have rights with respect to the party from which you acquired the software. These license terms do not change those other rights if the laws of your state or country do not permit it to do so. For example, if you acquired the software in one of the below regions, or if mandatory country law applies, then the following provisions apply to you:
Australia. You have statutory guarantees under the Australian Consumer Law and nothing in these license terms is intended to affect those rights.
Canada. You may stop receiving updates on your device by turning off Internet access. If and when you re-connect to the Internet, the software will resume checking for and installing updates. The product documentation, if any, may also specify how to turn off updates for your specific device or software.
Germany and Austria.
Warranty. The properly licensed software will perform substantially as described in any Microsoft materials that accompany the software. However, Microsoft gives no contractual guarantee in relation to the licensed software.
Limitation of Liability. In case of intentional conduct, gross negligence, claims based on the Product Liability Act, as well as, in case of death or personal or physical injury, Microsoft is liable according to the statutory law.
Subject to the foregoing clause (ii), Microsoft will only be liable for slight negligence if Microsoft is in breach of such material contractual obligations, the fulfillment of which facilitate the due performance of these license terms, the breach of which would endanger the purpose of these license terms and the compliance with which a party may constantly trust in (so-called "cardinal obligations"). In other cases of slight negligence, Microsoft will not be liable for slight negligence.
DISCLAIMER OF WARRANTY. The software is licensed “as-is.” You bear the risk of using it. Microsoft gives no express warranties, guarantees or conditions. To the extent permitted under your local laws, Microsoft excludes the implied warranties of merchantability, fitness for a particular purpose and non-infringement.
LIMITATION ON DAMAGES. You can recover from Microsoft and its suppliers only direct damages up to U.S. $5.00. You cannot recover any other damages, including consequential, lost profits, special, indirect or incidental damages.
This limitation applies to (a) anything related to the software, services, content (including code) on third party Internet sites, or third party applications; and (b) claims for breach of contract, breach of warranty, guarantee or condition, strict liability, negligence, or other tort to the extent permitted by applicable law.
It also applies even if Microsoft knew or should have known about the possibility of the damages. The above limitation or exclusion may not apply to you because your country may not allow the exclusion or limitation of incidental, consequential or other damages.

Please note: As this software is distributed in Canada, some of the clauses in this agreement are provided below in French. 
Remarque: Ce logiciel étant distribué au Canada, certaines des clauses dans ce contrat sont fournies ci-dessous en français. 
EXONÉRATION DE GARANTIE. Le logiciel visé par une licence est offert « tel quel ». Toute utilisation de ce logiciel est à votre seule risque et péril. Microsoft n’accorde aucune autre garantie expresse. Vous pouvez bénéficier de droits additionnels en vertu du droit local sur la protection des consommateurs, que ce contrat ne peut modifier. La ou elles sont permises par le droit locale, les garanties implicites de qualité marchande, d’adéquation à un usage particulier et d’absence de contrefaçon sont exclues. 
LIMITATION DES DOMMAGES-INTÉRÊTS ET EXCLUSION DE RESPONSABILITÉ POUR LES DOMMAGES. Vous pouvez obtenir de Microsoft et de ses fournisseurs une indemnisation en cas de dommages directs uniquement à hauteur de 5,00 $ US. Vous ne pouvez prétendre à aucune indemnisation pour les autres dommages, y compris les dommages spéciaux, indirects ou accessoires et pertes de bénéfices. 
Cette limitation concerne: 
tout ce qui est relié au logiciel, aux services ou au contenu (y compris le code) figurant sur des sites Internet tiers ou dans des programmes tiers; et 
les réclamations au titre de violation de contrat ou de garantie, ou au titre de responsabilité stricte, de négligence ou d’une autre faute dans la limite autorisée par la loi en vigueur. 
Elle s’applique également, même si Microsoft connaissait ou devrait connaître l’éventualité d’un tel dommage. Si votre pays n’autorise pas l’exclusion ou la limitation de responsabilité pour les dommages indirects, accessoires ou de quelque nature que ce soit, il se peut que la limitation ou l’exclusion ci-dessus ne s’appliquera pas à votre égard. 
EFFET JURIDIQUE. Le présent contrat décrit certains droits juridiques. Vous pourriez avoir d’autres droits prévus par les lois de votre pays. Le présent contrat ne modifie pas les droits que vous confèrent les lois de votre pays si celles-ci ne le permettent pas.

EULA ID: Cpp_v14_ENU.1033
