# A blockchain‑based decentralized efficient investigation framework for IoT digital forensics
Authors: Jung Hyun Ryu · Pradip Kumar Sharma · Jeong Hoon Jo · Jong Hyuk Park
Journal: The Journal of Supercomputing
DOI: 10.1007/s11227-019-02779-9

Private BlockChain based framework built on top of Ethereum. Decentralized data in different blockchain nodes. Chain of Custody on blockchain. Forensic data preservation through blockchain (distributed db with multiple authorities). More transparent investigation process. Anyone involved in the investigation process can verify integrity of data.

Limitations
- Possible breach of privacy, makes n/w data available for multiple stake-holders
- Only focuses on data collection and preservation

# A Framework for Digital Forensics Analysis based on Semantic Role Labeling
Authors: Ravi Barreira · Vládia Pinheiro · Vasco Furtado
Conference: IEEE International Conference on Intelligence and Security Informatics, ISI
DOI: 10.1109/ISI.2017.8004876

An AI(NLP based digital forensics framework based on Frame Semantics and FrameFOR knowledge base. The paper aims to solve the problem of volume of text & the amount of time it takes to manually read and understand it. The paper also tries to overcome limitations of traditional search tools that make use of keywords & their inherent need of vast database of keywords, by making use of Semantic Role Labeling, the paper tries to extract context out of text and their meaning, to increase analysis speed. It also addresses need for a uniform & standardized lexicon of keywords for weapons, intoxicants & other relevant identifiers. Main knowledege bases for SRL are VerbNet & FrameNet. FrameFOR, however, is built on specific frames, that were selected from FrameNet manually, by a cyber crime forensics expert to construct FrameFOR database. Aim of this paper includes identification of objects, people, entities in text and different terms being used by criminals to disguise their intent.

More info on VerbNet & FrameNet:
https://joanbanach.wordpress.com/2012/02/05/verbnet-and-framenet-lexical-semantics-iii/

Limitations:
 - Relies on manual updation of FrameFOR db
 - Removes emojis during routine cleaning (routines are run before pre-processing)

# A new network forensic framework based on deeplearning for Internet of Things networks: A particle deep framework
Authors: Nickolaos Koroniotis · Nour Moustafa · Elena Sitnikova
Journal: Future Generation Computer Systems
DOI: 10.1016/j.future.2020.03.042

A deep learning & optimization based n/w forensics framework that performs flow analysis on collected network data. The framework makes use of Particle Swarm Optimization(PSO) for tuning hyperparameters of the DNN model. The paper compares their framework with others on Bot-IOT & UNSW-NB15 datasets. The framework makes use of network flow data for analysis and avoids any changes within the existing IoT. Steps include capturing data, collecting data, sending it to the PSO algorithm for tuning hyper parameters of the DNN model, using the MLP DNN for classifying data into two categories and then measuring its performance.

- Figuring out spoofed IPs is challenging

# A multilayered semantic framework for integrated forensic acquisitionon social media
Authors: Humaira Arshad · Aman Jantan · Gan Keng Hoon · Anila Sahar Butt
Journal: Digital Investigation
DOI: 10.1016/j.diin.2019.04.002

A digital forensics framework for investigating Online Social Networks(OSN). The goal of this framework is to provide a semi-automated tool for data acquisition, analysis and visualisation. The framework understands the hetrogenity of data among different OSNs and have come up with a hybrid ontology approach to save and correlate data, so that a timeline can be established among different events. The framework makes use of officially available public APIs of popular OSNs like Facebook, Twitter to archive data.

Limiations:
- The framework does not employ machine learning or artificial intelligence for extracting sentiment, context, objects from social media posts
- While the paper talks aobut how cyber criminals make use of social media, the framework does not talk about the possibility of malicious files, stegnographic content and other similar data or encoded text that may require further manual analysis

# A machine learning-based FinTech cyber threat attribution framework using high-level indicators of compromise
Authors: Umara Noor · Zahid Anwar · Tehmina Amjad · Kim-Kwang Raymond Choo
Journal: Future Generation Computer Systems
DOI: 10.1016/j.future.2019.02.013

An ML based framework for cyber threat attribution that makes use of unstructures CTI reports and extracts relevant high level indicators of compromise(IoC). The framework has three phases, collecting CTAs/APTs & CTI documents, mapping TTPs query labels with unstructured CTI & finally CTA class prediciton. The intial dataset is prepared from ATT&CK taxonomy, provided by MITRE. The framework makes use of Latent Semantic Analysis to index CTI reports. The framework trains classifiers on a correlation matrix of CTAs & TTPs to automate investigative process and attribute the cyber attacks, reports suggest that the automated classifers achieve upto 94% accuracy as compared to manually investigated IoCs.

Limitations
- The framework is dependent on availability of correct threat data

# Towards a practical cloud forensics logging framework

A cloud logging framework for performing cloud fornesic tasks. The framework is embedded into the hypervysor, it generates one or more log files per Cloud Service User(CSU). All the files generated by this framework are stored in a pred-defined location. The framework conforms to ACPO guidelines, validated against NIST draft report on Cloud Forensic science challenges.

Limitations:
- The framework needs to be installed and embedded directly into the CSPs hypervysor